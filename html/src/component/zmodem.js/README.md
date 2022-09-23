# zmodem.js - ZMODEM for JavaScript

[![build status](https://api.travis-ci.org/FGasper/zmodemjs.svg?branch=master)](http://travis-ci.org/FGasper/zmodemjs)

# SYNOPSIS

    let zsentry = new Zmodem.Sentry( {
        to_terminal(octets) { .. },  //i.e. send to the terminal

        sender(octets) { .. },  //i.e. send to the ZMODEM peer

        on_detect(detection) { .. },  //for when Sentry detects a new ZMODEM

        on_retract() { .. },  //for when Sentry retracts a Detection
    } );

    //We have to configure whatever gives us new input to send that
    //input to zsentry.
    //
    //EXAMPLE: From web browsers that use WebSocket …
    //
    ws.addEventListener("message", function(evt) {
        zsentry.consume(evt.data);
    } );

The `on_detect(detection)` function call is probably the most complex
piece of the above; one potential implementation might look like:

    on_detect(detection) {

        //Do this if we determine that what looked like a ZMODEM session
        //is actually not meant to be ZMODEM.
        if (no_good) {
            detection.deny();
            return;
        }

        zsession = detection.confirm();

        if (zsession.type === "send") {

            //Send a group of files, e.g., from an <input>’s “.files”.
            //There are events you can listen for here as well,
            //e.g., to update a progress meter.
            Zmodem.Browser.send_files( zsession, files_obj );
        }
        else {
            zsession.on("offer", (xfer) => {

                //Do this if you don’t want the offered file.
                if (no_good) {
                    xfer.skip();
                    return;
                }

                xfer.accept().then( () => {

                    //Now you need some mechanism to save the file.
                    //An example of how you can do this in a browser:
                    Zmodem.Browser.save_to_disk(
                        xfer.get_payloads(),
                        xfer.get_details().name
                    );
                } );
            });

            zsession.start();
        }
    }

# DESCRIPTION

zmodem.js is a JavaScript implementation of the ZMODEM
file transfer protocol, which facilitates file transfers via a terminal.

# STATUS

This library is BETA quality. It should be safe for general use, but
breaking changes may still happen.

# HOW TO USE THIS LIBRARY

The basic workflow is:

1. Create a `Zmodem.Sentry` object. This object must scan all input for
a ZMODEM initialization string. See `zsentry.js`’s documentation for more
details.

2. Once that initialization is found, the `on_detect` event is fired
with a `Detection` object as parameter. At this point you can `deny()`
that Detection or `confirm()` it; the latter will return a Session
object.

3. Now you do the actual file transfer(s):

    * If the session is a receive session, do something like this:

            zsession.on("offer", (offer) => { ... });
                let { name, size, mtime, mode, serial, files_remaining, bytes_remaining } = offer.get_details();

                offer.skip();

                //...or:

                offer.on("input", (octets) => { ... });

                //accept()’s return resolves when the transfer is complete.
                offer.accept().then(() => { ... });
            });
            zsession.on("session_end", () => { ... });
            zsession.start();

        The `offer` handler receives an Offer object. This object exposes the details
    about the transfer offer. The object also exposes controls for skipping or
    accepting the offer.

    * Otherwise, your session is a send session. Now the user chooses
zero or more files to send. For each of these you should do:

            zsession.send_offer( { ... } ).then( (xfer) => {
                if (!xfer) ... //skipped

                else {
                    xfer.send( chunk );
                    xfer.end( chunk ).then(after_end);
                }
            } );

        Note that `xfer.end()`’s return is a Promise. The resolution of this
Promise is the point at which either to send another offer or to do:

            zsession.close().then( () => { ... } );

        The `close()` Promise’s resolution is the point at which the session
has ended successfully.

That should be all you need. If you want to go deeper, though, each module
in this distribution has JSDoc and unit tests.

# RATIONALE

ZMODEM facilitates terminal-based file transfers.
This was an important capability in the 1980s and early 1990s because
most modem use was for terminal applications, especially
[BBS](https://en.wikipedia.org/wiki/Bulletin_board_system)es.
(This was how, for example,
popular shareware games like [Wolfenstein 3D](http://3d.wolfenstein.com)
were often distributed.) The World Wide Web in the
mid-1990s, however, proved a more convenient way to accomplish most of
what BBSes were useful for, as a result of which the problem that ZMODEM
solved became a much less important one.

ZMODEM stuck around, though, as it remained a convenient solution
for terminal users who didn’t want open a separate session to transfer a
file. [Uwe Ohse](https://uwe.ohse.de/)’s
[lrzsz](https://ohse.de/uwe/software/lrzsz.html) package
provided a portable C implementation of the protocol (reworked from
the last public domain release of the original code) that is installed on
many systems today.

Where `lrzsz` can’t reach, though, is terminals that don’t have command-line
access—such as terminals that run in JavaScript. Now that
[WebSocket](https://en.wikipedia.org/wiki/WebSocket) makes real-time
applications like terminals possible in a web browser,
there is a use case for a JavaScript
implementation of ZMODEM to allow file transfers in this context.

# GENERAL FLOW OF A ZMODEM SESSION:

The following is an overview of an error-free ZMODEM session.

0. If you call the `sz` command (or equivalent), that command will send
a special ZRQINIT “pre-header” to signal your terminal to be a ZMODEM
receiver.

1. The receiver, upon recognizing the ZRQINIT header, responds with
a ZRINIT header.

2. The sender sends a ZFILE header along with information about the file.
(This may also include the size and file count for the entire batch of files.)

3. The recipient either accepts the file or skips it.

4. If the recipient did not skip the file, then the sender sends the file
contents. At the end the sender sends a ZEOF header to let the recipient
know this file is done.

5. The recipient sends another ZRINIT header. This lets the sender know that
the recipient confirms receipt of the entire file.

6. Repeat steps 2-5 until the sender has no more files to send.

7. Once the sender has no more files to send, the sender sends a ZEOF header,
which the recipient echoes back. The sender closes the session by sending
`OO` (“over and out”).

# PROTOCOL NOTES AND ASSUMPTIONS

Here are some notes about this particular implementation.

Particular notes:

* We send with a maximum data subpacket size of 8 KiB (8,192 bytes). While
the ZMODEM specification stipulates a maximum of 1 KiB, `lrzsz` accepts
the larger size, and it seems to have become a de facto standard extension
to the protocol.

* Remote command execution (i.e., ZCOMMAND) is unimplemented. It probably
wouldn’t work in browsers, which is zmodem.js’s principal use case.

* No file translations are done. (Unix/Windows line endings are a
future feature possibility.)

* It is assumed that no error correction will be needed. All connections
are assumed to be **“reliable”**; i.e.,
data is received exactly as sent. We take this for granted today,
but ZMODEM’s original application was over raw modem connections that
often didn’t have reliable hardware error correction. TCP also wasn’t
in play to do software error correction as generally happens
today over remote connections. Because the forseeable use of zmodem.js
is either over TCP or a local socket—both of which are reliable—it seems
safe to assume that zmodem.js will not need to implement error correction.

* zmodem.js sends with CRC-16 by default. Ideally we would just use CRC-16
for everything, but lsz 0.12.20 has a [buffer overflow bug](https://github.com/gooselinux/lrzsz/blob/master/lrzsz-0.12.20.patch) that rears its
head when you try to abort a ZMODEM session in the middle of a CRC-16 file
transfer. To avoid this bug, zmodem.js advertises CRC-32 support when it
receives a file, which makes lsz avoid the buffer overflow bug by using
CRC-32.

    The bug is reported, incidentally, and a fix is expected (nearly 20 years
    after the last official lrzsz release!).

* There is no XMODEM/YMODEM fallback.

* Occasionally lrzsz will output things to the console that aren’t
actual ZMODEM—for example, if you skip an offered file, `sz` will write a
message about it to the console. For the most part we can accommodate these
because they happen between ZMODEM headers; however, it’s possible to
“poison” such messages, e.g., by sending a file whose name includes a
ZMODEM header. So don’t do that. :-P

# IMPLEMENTATION NOTES

* I initially had success integrating zmodem.js with
[xterm.js](https://xtermjs.org); however, that library’s plugin interface
changed dramatically, and I haven’t created a new plugin to replace the
old one. (It should be relatively straightforward if someone else wants to
pick it up.)

* Browsers don’t have an easy way to download only part of a file;
as a result, anything the browser saves to disk must be the entire file.

* ZMODEM is a _binary_ protocol. (There was an extension planned
to escape everything down to 7-bit ASCII, but it doesn’t seem to have
been implemented?) Hence, **if you use WebSocket, you’ll need to use
binary messages, not text**.

* lrzsz is the only widely-distributed ZMODEM implementation nowadays,
which makes it a de facto standard in its
own right. Thus far all end-to-end testing has been against it. It is
thus possible that resolutions to disparities between `lrzsz` and the
protocol specification may need to favor the implementation.

* It is a generally-unavoidable byproduct of how ZMODEM works that
the first header in a ZMODEM session will echo to the terminal. This
explains the unsightly `**B0000…` stuff that you’ll see when you run
either `rz` or `sz`.

    That header
    will include some form of line break. (From `lrzsz` means bytes 0x0d
    and 0x8a—**not** 0x0a). Your terminal might react oddly to that;
    if it does, try stripping out one or the other line ending character.

# PROTOCOL CHOICE

Both XMODEM and YMODEM (including the latter’s many variants) require the
receiver to initiate the session by sending a “magic character” (ASCII SOH);
the problem is that there’s nothing in the protocol to prompt the receiver
to do so. ZMODEM is sender-driven, so the terminal can show a notice that
says, “Do you want to receive a file?”

This is a shame because these other two protocols are a good deal simpler
than ZMODEM. The YMODEM-g variant in particular would be better-suited to
our purpose because it doesn’t “litter” the transfer with CRCs.

There is also [Kermit](http://www.columbia.edu/kermit/kermit.html), which
seems to be more standardized than ZMODEM but **much** more complex.

# DESIGN NOTES

zmodem.js tries to avoid “useless” states:
either we fail completely, or we succeed. To that end, some callbacks are
required arguments (e.g., the Sentry constructor’s `to_terminal` argument),
while others are registered separately.

Likewise, for this reason some of the session-level logic is exposed only
through the Transfer and Offer objects. The Session creates these
internally then exposes them via callback

# SOURCES

ZMODEM is not standardized in a nice, clean, official RFC like DNS or HTTP;
rather, it was one guy’s solution to a particular problem. There is
documentation, but it’s not as helpful as it might be; for example,
there’s only one example workflow given, and it’s a “happy-path”
transmission of a single file.

As part of writing zmodem.js I’ve culled together various resources
about the protocol. As far as I know these are the best sources for
information on ZMODEM.

Two documents that describe ZMODEM are saved in the repository for reference.
The first is the closest there is to an official ZMODEM specification:
a description of the protocol from its author, Chuck Forsberg. The second
seems to be based on the first and comes from
[Jacques Mattheij](https://jacquesmattheij.com).

**HISTORICAL:** The included `rzsz.zip` file (fetched from [ftp://archives.thebbs.org/file_transfer_protocols/](ftp://archives.thebbs.org/file_transfer_protocols/) on 16 October 2017)
is the last public domain release
from Forsberg. [http://freeware.nekochan.net/source/rzsz/](http://freeware.nekochan.net/source/rzsz/) has what is supposedly Forsberg’s last shareware release;
I have not looked at it except for the README. I’m not sure of the
copyright status of this software: Forsberg is deceased, and his company
appears to be defunct. Regardless, neither it nor its public domain
predecessor is likely in widespread use.

Here are some other available ZMODEM implementations:

* [lrzsz](https://ohse.de/uwe/software/lrzsz.html)

    A widely-deployed adaptation of Forsberg’s last public domain ZMODEM
    code. This is the de facto “reference” implementation, both by virtue
    of its wide availability and its derivation from Forsberg’s original.
    If your server has the `rz` and `sz` commands, they’re probably
    from this package.

* [SyncTERM](http://syncterm.bbsdev.net)

    Based on Jacques Mattheij’s ZMODEM implementation, originally called
    zmtx/zmrx. This is a much more readable implementation than lrzsz
    but lamentably one that doesn’t seem to compile as readily.

* [Qodem](https://github.com/klamonte/qodem)

    This terminal emulator package appears to contain its own ZMODEM
    implementation.

* [PD Zmodem](http://pcmicro.com/netfoss/pdzmodem.html)

    I know nothing of this one.

* [zmodem (Rust)](https://github.com/lexxvir/zmodem)

    A pure [Rust](http://rust-lang.org) implementation of ZMODEM.

# REQUIREMENTS

This library only supports modern browsers. There is no support for
Internet Explorer or other older browsers planned.

The tests have run successfully against node.js version 8.

# DOCUMENTATION

Besides this document, each module has inline [jsdoc](http://usejsdoc.org).
You can see it by running `yarn` in the repository’s root directory;
the documentation will build in a newly-created `documentation` directory.

# CONTRIBUTING

Contributions are welcome via
[https://github.com/FGasper/zmodemjs](https://github.com/FGasper/zmodemjs).

# TROUBLESHOOTING

Before you do anything else, set `Zmodem.DEBUG` to true. This will log
useful information about the ZMODEM session to your JavaScript console. That
may give you all you need to fix your problem.

If you have trouble transferring files, try these diagnostics:

1. Transfer an empty file. (Run `touch empty.bin` to create one named `empty.bin`.)

2. Transfer a small file. (`echo hello > small.txt`)

3. Transfer a file that contains all possible octets. (`perl -e 'print chr for 0 .. 255' > all_bytes.bin`)

4. If a specific file fails, does it still fail if you `truncate` a copy of
the file down to, say, half size and transfer that truncated file? Does it
work if you truncate the file down to 1 byte? If so, then use this method
to determine which specific place in the file triggers the transfer error.

**IF YOU HAVE DONE THE ABOVE** and still think the problem is with zmodem.js,
you can file a bug report. Note that, historically, most bug reports have
reflected implementation errors rather than bugs in zmodem.js.

# TODO

* Teach send sessions to “fast-forward” so as to honor requests for
append-style sessions.

* Implement newline conversions.

* Teach Session how to do and to handle pre-CRC checks.

* Possible: command-line `rz`, if there’s demand for it, e.g., in
environments where `lrzsz` can’t run. (NB: The distribution includes
a bare-bones, proof-of-concept `sz` replacement.)

# KNOWN ISSUERS

* In testing, Microsoft Edge appeared not to care what string was given
to `<a>`’s `download` attribute; the saved filename was based on the
browser’s internal Blob object URL instead.

# COPYRIGHT

Copyright 2017 Gasper Software Consulting

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Parts of the CRC-16 logic are adapted from crc-js by Johannes Rudolph.
