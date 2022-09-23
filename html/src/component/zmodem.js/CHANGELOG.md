# 0.1.10

Make unrecognized-header detection more resilient.

Ignore extra ZRPOS if received while sending a file. (See comments
for the rationale.)

Expose Zmodem.DEBUG for runtime adjustment.

Add a proof-of-concept CLI “sz” implementation to the distribution.

Change quality designation from ALPHA to BETA.

Documentation updates, including addition of a TROUBLESHOOTING section.

---

# 0.1.9

No production changes; this just disables a flapping test.

---

# 0.1.8

This version introduces some minor, mostly-under-the-hood changes:

1. `accept()` callbacks now fire after receipt of the ZEOF.
Previously they didn’t fire until the sender indicated either the next
file (ZFILE) or the end of the batch (ZFIN). This actually brings the
behavior more in line with the documentation.

2. In the same vein, the `file_end` event now fires before ZRINIT is sent.

3. `skip()` is now a no-op if called outside of a transfer. Previously
it always sent a ZSKIP, which confused `sz` into sending an extra ZFIN if
it happened outside of a file transfer, which tripped up protocol errors
in zmodem.js.

4. A misnamed variable is now fixed.

Additionally, a bug in the tests that caused the test runner to skip
some test files is fixed. Every test now runs, and new tests are added that
verify the “happy-path” in receive sessions.
