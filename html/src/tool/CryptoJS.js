/*
* 账号加密
*/
import CryptoJS from 'crypto-js';

export default {
    //加密
    encrypt(word, key) {
        const keyParse = CryptoJS.enc.Utf8.parse(key)
        const srcs = CryptoJS.enc.Utf8.parse(word)
        const encrypted = CryptoJS.AES.encrypt(srcs, keyParse, {
            iv: keyParse,
            mode: CryptoJS.mode.CBC,
            padding: CryptoJS.pad.Pkcs7
        })
        return encrypted.toString()
    },
    //解密
    decrypt(word, key) {
        const keyParse = CryptoJS.enc.Utf8.parse(key)
        const decrypted = CryptoJS.AES.decrypt(word, keyParse, {
            iv: keyParse,
            mode: CryptoJS.mode.CBC,
            padding: CryptoJS.pad.Pkcs7
        })
        return CryptoJS.enc.Utf8.stringify(decrypted).toString()
    },
    //随机生成指定数量的16进制key
    generatekey(num) {
        let library = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
        let key = "";
        for (var i = 0; i < num; i++) {
            let randomPoz = Math.floor(Math.random() * library.length);
            key += library.substring(randomPoz, randomPoz + 1);
        }
        return key;
    }
}