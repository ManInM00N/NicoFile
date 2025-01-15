<script setup>
import { RouterLink, RouterView } from 'vue-router'
import axios from 'axios'
import {ref} from "vue";
import {getFileMd5, sleep} from "@/assets/js/encyrpt.js";

const progress = ref(0),progressNow = ref('')
const filename = ref('')
const hasupload = ref(0)
const chunktotals = ref(0)
const Ext = ref('')
const size = ref(0)
const fileMd5 = ref('')
const maxSize  = ref(5 * 1024 * 1024 * 1024), // 上传最大文件限制  最小单位是b
    chunkSize = 1024 * 1024 * 5, // 每块文件大小   100mb
    fileList = ref([])

async function Remove(file, filelist) {
  console.log(file.name,file.status)
  fileList.value = filelist
  hasupload.value = 0
  size.value = 0
  progress.value = 0
  progressNow.value = ''
}
async function CalcFile(file, filelist) {
  fileList.value = filelist
  console.log(file.status)
  if (file.status === 'removed'){
    hasupload.value = 0
    size.value = 0
    progress.value = 0
    progressNow.value = ''
    return
  }
  let fileName = file.name
  Ext.value = ''
  if (fileName.lastIndexOf('.') !== -1) {
    Ext.value = fileName.substring(fileName.lastIndexOf('.'))
    filename.value = fileName.substring(0, fileName.lastIndexOf('.'))
  }
  console.log(filename,Ext)
  size.value = file.size
  chunktotals.value = Math.ceil(size.value / chunkSize);
  hasupload.value = 0
  fileMd5.value = await getFileMd5(file.raw)
  console.log(fileMd5.value)
}

const uploadFileToServer = async (file, chunkNumber, fileName,_md5) => {
  let form = new FormData();
  form.append("chunk", file);
  form.append("chunkIndex", chunkNumber);
  form.append("md5", _md5);
  form.append("filename", fileName);
  console.log(_md5,chunkNumber)
  const result = await axios.post("http://localhost:8888/api/v1/file/uploadchunk", form,{
    headers: {
      "Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEsImV4cCI6MTczNzMxNDk5OSwiaWF0IjoxNzM2NzExMzk5fQ.6ASy3-He6IxqhXmATyKekvGWtOw5I9PPb1_9-rgJNDs",
    },
    onUploadProgress: (progressEvent) => {
      if (progressEvent.lengthComputable) {
        console.log(progressEvent.loaded,progressEvent)
        const percent =  ((hasupload.value +progressEvent.loaded)  / 1024 / 1024).toFixed(2);
        const tot = (size.value/1024/1024).toFixed(2)
        progressNow.value = `${percent} MB/ ${tot} MB`
        progress.value = (hasupload.value +progressEvent.loaded) / size.value * 100
      }
    },
  })
  hasupload.value += file.size
  return result
}

const mergeFiles = async (chunkTotal, fileName,ext,size) => {
  const result =await axios.post("http://localhost:8888/api/v1/file/mergechunk", {
    chunkNum: chunkTotal,
    filename: fileName,
    md5 :fileMd5.value,
    ext: ext,
    size: size
  },{
    headers: {
      "Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEsImV4cCI6MTczNzMxNDk5OSwiaWF0IjoxNzM2NzExMzk5fQ.6ASy3-He6IxqhXmATyKekvGWtOw5I9PPb1_9-rgJNDs",
    }
  })
  return result
}

async function checkchunk(chunkTotal, fileName,md5arr) {
  const resp = await axios.post("http://localhost:8888/api/v1/file/checkchunk", {
    chunkNum: chunkTotal,
    filename: fileName,
    md5 : md5arr,
    fileMd5 : fileMd5.value,
    ext : Ext.value
  },{
    headers: {
      "Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEsImV4cCI6MTczNzMxNDk5OSwiaWF0IjoxNzM2NzExMzk5fQ.6ASy3-He6IxqhXmATyKekvGWtOw5I9PPb1_9-rgJNDs",
    }
  })
  return resp.data
}
const merge = async ()=>{
  const res = await mergeFiles(chunktotals.value, filename.value,Ext.value,size.value)
  console.log(res.data)
}
const submit = async () => {
  let uploadFile = fileList.value[0]
  if (chunktotals.value > 0) {
    let md5arr = []
    for (let chunkNumber = 0, start = 0; chunkNumber < chunktotals.value; chunkNumber++, start += chunkSize) {
      let end = Math.min(size.value, start + chunkSize);
      const file = uploadFile.raw.slice(start, end)
      const _md5  = await   getFileMd5(file)
      md5arr.push(_md5)
    }
    // 检查文件切片上传进度
    const resp = await checkchunk(chunktotals.value, filename.value,md5arr)
    console.log("Total Chunks: ",chunktotals.value,resp)
    hasupload.value = resp.accept * chunkSize
    // 继续上传
    for (let chunkNumber = resp.accept, start = resp.accept*chunkSize; chunkNumber < chunktotals.value; chunkNumber++, start += chunkSize) {
      let end = Math.min(size.value, start + chunkSize);
      const files = uploadFile.raw.slice(start, end)
      console.log(files,start,end)
      const result = await uploadFileToServer(files, chunkNumber , filename.value,md5arr[chunkNumber])
      if (result.status !== 200) {
        console.log(result.data,result)
        progress.value = 0
        progressNow.value = ''
        console.log("上传失败")
        return
      }
      await sleep(1000)
    }
    progress.value = 100
    progressNow.value = "上传完成"
  }
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
        :on-change="CalcFile"
        :on-remove="Remove"
    >
    </el-upload>
    <el-button size="small" type="primary" @click="submit">点击上传</el-button>
    <el-button size="small" type="primary" @click="merge">点击合并</el-button>
    <el-progress
        :text-inside="true"
        :percentage="progress"
        striped
        striped-flow
        :stroke-width="20"
        style="width:600px"
    >
        {{progressNow}}
    </el-progress>
  </div>
</template>

<style scoped>
</style>
