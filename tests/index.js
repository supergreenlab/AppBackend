const fs = require('fs')
const axios = require('axios')
const FormData = require('form-data')

const SERVER = 'http://localhost:8080'
const MINIO_SERVER = 'http://minio:9000'

const main = async() => {
  const { headers: { 'x-sgl-token': userJwt } } = await axios.post(`${SERVER}/login`, {
    handle: 'stant',
    password: 'stant',
  })

  const { headers: { 'x-sgl-token': jwt } } = await axios.post(`${SERVER}/userend`, {}, {
    headers: {
      'Authentication': `Bearer ${userJwt}`,
    },
  })

  const { data: { filePath, thumbnailPath } } = await axios.post(`${SERVER}/feedMediaUploadURL`, {
    fileName: 'logo.jpg',
  }, {
    headers: {
      'Authentication': `Bearer ${jwt}`,
    },
  })

  const stats = fs.statSync('logo.jpg');
  try {
    const resp = await axios.put(`${MINIO_SERVER}${filePath}`, fs.createReadStream('logo.jpg'), {
      headers: {
        'Content-Type': 'image/jpg',
        'Content-Length': stats['size'],
      }
    })
    console.log('pouet', resp)
  } catch(e) {
    console.log(e)
  }
}

main()
