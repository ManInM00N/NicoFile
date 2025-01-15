import SparkMD5 from "spark-md5";

export const sleep = ms => new Promise(r => setTimeout(r, ms));

export function getFileMd5(file) {
    return new Promise((resolve, reject) => {
        const fileReader = new FileReader();
        const Spark = new SparkMD5.ArrayBuffer();
        fileReader.onload = function (e) {
            try {
                Spark.append(e.target.result);
                const md5Value = Spark.end(); // 获取最终的 MD5 值
                resolve(md5Value);
            } catch (error) {
                reject(error);
            }
        };

        fileReader.onerror = function () {
            reject(new Error("FileReader error occurred"));
        };

        fileReader.readAsArrayBuffer(file); // 读取文件为 ArrayBuffer
    });
}