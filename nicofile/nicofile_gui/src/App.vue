<script setup>
import { RouterLink, RouterView } from 'vue-router'
import axios from 'axios'
import SparkMD5 from 'spark-md5'
import {ref} from "vue";
let maxSize  = ref(5 * 1024 * 1024 * 1024), // 上传最大文件限制  最小单位是b
    chunkSize = 1024 * 1024 * 5, // 每块文件大小   100mb
    fileList = ref([])

function  loadJsonFromFile(file,filelist) {
  fileList.value = filelist
  console.log(fileList,fileList.value[0])
}
const uploadFileToServer = async (file, chunkNumber, fileName) => {
  let Spark=new SparkMD5.ArrayBuffer()
  // Spark.append(file)
  // let _md5 = Spark.end(false)
  var  fileReader=new FileReader()
  const _md5 = await new Promise((resolve, reject) => {
    fileReader.onload = function (e) {
      try {
        Spark.append(e.target.result);
        resolve(Spark.end());
      } catch (error) {
        reject(error);
      }
    };
    fileReader.onerror = function () {
      reject(new Error("FileReader error occurred"));
    };
    fileReader.readAsArrayBuffer(file);
  });

  let form = new FormData();
  form.append("chunk", file);
  form.append("chunkIndex", chunkNumber);
  form.append("md5", _md5);
  form.append("filename", fileName);
  console.log(_md5,chunkNumber)
  // let result = {data:''}
  const result = await axios.post("http://localhost:8888/api/v1/file/uploadchunk", form,{
    headers: {
      "Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEsImV4cCI6MTczNzMxNDk5OSwiaWF0IjoxNzM2NzExMzk5fQ.6ASy3-He6IxqhXmATyKekvGWtOw5I9PPb1_9-rgJNDs",
    }
  })
  return result
}
const mergeFiles = async (chunkTotal, fileName,ext,size) => {
  const result =await axios.post("http://localhost:8888/api/v1/file/mergechunk", {
    chunkNum: chunkTotal,
    filename: fileName,
    md5 :"",
    ext: ext,
    size: size
  },{
    headers: {
      "Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEsImV4cCI6MTczNzMxNDk5OSwiaWF0IjoxNzM2NzExMzk5fQ.6ASy3-He6IxqhXmATyKekvGWtOw5I9PPb1_9-rgJNDs",
    }
  })
  return result.data
}
const filename = ref('')
const chunktotals = ref(0)
const Ext = ref('')
const size = ref(0)
const merge = async ()=>{
  const res = await mergeFiles(chunktotals.value, filename.value,Ext.value,size.value)
  console.log(res.data)
}
const submit = async () => {
  let uploadFile = fileList.value[0]
  let fileName = uploadFile.name
  let ext = ""
  if (uploadFile.name.lastIndexOf('.') !== -1) {
    ext = uploadFile.name.substring(uploadFile.name.lastIndexOf('.'))
    fileName = uploadFile.name.substring(0, uploadFile.name.lastIndexOf('.'))
  }
  console.log(fileName,ext)
  const fileSize = uploadFile.size || 0
  let chunkTotals = Math.ceil(fileSize / chunkSize);
  console.log(chunkTotals)
  if (chunkTotals > 0) {
    for (let chunkNumber = 0, start = 0; chunkNumber < chunkTotals; chunkNumber++, start += chunkSize) {
      let end = Math.min(fileSize, start + chunkSize);
      const files = uploadFile.raw.slice(start, end)
      console.log(files,start,end)
      const result = await uploadFileToServer(files, chunkNumber , fileName)
      console.log(result.data)
    }
  }
  Ext.value  = ext
  chunktotals.value = chunkTotals
  size.value = fileSize
  filename.value = fileName

}
</script>

<template>
  <header>


  </header>
  <div >
    <el-upload
        ref = "upload"
        drag
        action="#"
        :auto-upload="false"
        :file-list="fileList"
        :limit="1"
        :on-change="loadJsonFromFile"
    >
    </el-upload>
    <el-button size="small" type="primary" @click="submit">点击上传</el-button>
    <el-button size="small" type="primary" @click="merge">点击合并</el-button>

  </div>
</template>

<style scoped>
</style>
