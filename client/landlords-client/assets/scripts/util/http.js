const httpUrl = 'http://172.21.165.80:80'
/**
 * Http 请求封装
 */
 const http = {
    login: '/login',
    register: '/register',
    logout: '/logout',
    /**
     * get请求
     * @param {string} url 
     * @param {function} callback 
     */
    get(url, callback) {
        let xhr = cc.loader.getXMLHttpRequest();
        xhr.onreadystatechange = function () {
            // cc.log('url', httpUrl, 'Get: readyState=' + xhr.readyState + '  xhr.status=' + xhr.status);
            if (xhr.readyState === 4 && xhr.status == 200) {
                let respone = xhr.responseText;
                let rsp = JSON.parse(respone);
                callback(rsp);
            }
        };
        xhr.withCredentials = true;
        xhr.open('GET', httpUrl + url, true);

        // if (cc.sys.isNative) {
        // xhr.setRequestHeader('Access-Control-Allow-Origin', '*');
        // xhr.setRequestHeader('Access-Control-Allow-Methods', 'GET, POST');
        // xhr.setRequestHeader('Access-Control-Allow-Headers', 'x-requested-with,content-type,authorization');
        // xhr.setRequestHeader("Content-Type", "application/json");
        // xhr.setRequestHeader('Authorization', 'Bearer ' + cc.myGame.gameManager.getToken());
        // xhr.setRequestHeader('Authorization', 'Bearer ' + "");
        // }

        // note: In Internet Explorer, the timeout property may be set only after calling the open()
        // method and before calling the send() method.
        xhr.timeout = 8000;// 8 seconds for timeout

        xhr.send();
    },

    /**
     * post请求
     * @param {string} url 
     * @param {object} params 
     * @param {function} callback 
     */
    post(url, params, callback) {
        let xhr = cc.loader.getXMLHttpRequest();
        xhr.onreadystatechange = function () {
            // cc.log('url', httpUrl, 'xhr.readyState=' + xhr.readyState + '  xhr.status=' + xhr.status);
            if (xhr.readyState === 4 && xhr.status == 200) {
                let respone = xhr.responseText;
                let rsp = JSON.parse(respone);
                callback(rsp);
            }
        };
        xhr.withCredentials = true;
        xhr.open('POST', httpUrl + url, true);
        // if (cc.sys.isNative) {
        // xhr.setRequestHeader('Access-Control-Allow-Origin', '*');
        // xhr.setRequestHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
        // xhr.setRequestHeader('Access-Control-Allow-Headers', 'x-requested-with,content-type');
        xhr.setRequestHeader("Content-Type", "application/json");
        // xhr.setRequestHeader('Authorization', 'Bearer ' + cc.myGame.gameManager.getToken());
        // }

        // note: In Internet Explorer, the timeout property may be set only after calling the open()
        // method and before calling the send() method.
        xhr.timeout = 8000;// 8 seconds for timeout

        xhr.send(JSON.stringify(params));
    },
}

module.exports = http
