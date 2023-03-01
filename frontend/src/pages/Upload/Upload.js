import { useState, useEffect } from 'react'
import Box from '@mui/material/Box'
import Typography from '@mui/material/Typography'
import FileUpload from 'react-material-file-upload'
import { ReactNotifications } from 'react-notifications-component'
import 'react-notifications-component/dist/theme.css'
import { Button } from '@mui/material'
import CachedIcon from '@mui/icons-material/Cached'

import Header from '../../components/Header'
import Toast from '../../components/Toast'
import CollapsibleTable from '../../components/CollapsibleTable'
import * as AffiliateAPI from '../../services/affiliate-api'

export const Upload = () => {
  const [files, setFiles] = useState([])
  const [items, setItems] = useState([])

  async function fetchData() {
    try {
      const { data } = await AffiliateAPI.GetFiles()
      setItems(data)
    } catch ({ response }) {
      let err
      if (response && response.status > 399) {
        if (response.data.errors) {
          err = response.data.errors[0]
        } else {
          err = response.data.message
        }
      } else {
        err = 'An unknown error has occurred'
      }
      Toast({
        type: 'danger',
        title: 'Oops!',
        message: `Error: ${err}`,
      })
    }
  }

  useEffect(() => {
    fetchData()
  }, [files])

  const saveFile = async (files) => {
    if (files.length > 1) {
      return Toast({
        type: 'danger',
        title: 'Oops!',
        message: 'Não é possível realizar upload de vários arquivos!',
      })
    }

    try {
      const { data } = await AffiliateAPI.SendFile(files[0])
      Toast({
        type: 'success',
        title: 'Success!',
        message: `Arquivo enviado! ID: ${data.batchID}`,
      })
    } catch ({ response }) {
      if (response.status > 399) {
        const [err] = response.data.errors
        Toast({
          type: 'danger',
          title: 'Oops!',
          message: `Erro ao enviar o arquivo: ${err}`,
        })
      }
    } finally {
      setFiles([])
    }
  }

  return (
    <>
      <Header />
      <ReactNotifications />
      <Box
        sx={{
          marginTop: 8,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
        <Typography
          variant="h5"
          align="center"
          color="text.secondary"
          paragraph
        >
          Sales file upload
        </Typography>
        <Box component="form" noValidate sx={{ mt: 1, width: '75%' }}>
          <FileUpload value={files} onChange={saveFile} />
        </Box>

        <Box
          sx={{
            marginTop: 8,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            width: '75%',
          }}
        >
          <Box
            sx={{
              display: 'flex',
              flexDirection: 'row',
              alignItems: 'flex-end',
              justifyContent: 'flex-end',
              width: '100%',
              mb: 2,
            }}
          >
            <Button
              variant="outlined"
              startIcon={<CachedIcon />}
              onClick={fetchData}
            >
              Refresh
            </Button>
          </Box>

          <CollapsibleTable items={items} />
        </Box>
      </Box>
    </>
  )
}
