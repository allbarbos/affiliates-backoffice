import {
  Box,
  Divider,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
} from '@mui/material'
import { format } from 'date-fns/esm'

export const SpanningTable = ({ items, summary }) => {
  return (
    <TableContainer component={Paper}>
      <Table sx={{ minWidth: 700 }} aria-label="spanning table">
        <TableHead>
          <TableRow>
            <TableCell align="left">Product</TableCell>
            <TableCell align="left">Value</TableCell>
            <TableCell align="left">Seller</TableCell>
            <TableCell align="left">Type</TableCell>
            <TableCell align="left">Date</TableCell>
            <TableCell align="left">Batch ID</TableCell>
          </TableRow>
        </TableHead>
        <TableBody sx={{ position: 'relative' }}>
          {items.map((item, idx) => {
            return (
              <TableRow key={idx}>
                <TableCell align="left">{item.product}</TableCell>
                <TableCell align="left">
                  {item.value.toLocaleString('pt-br', {
                    style: 'currency',
                    currency: 'BRL',
                  })}
                </TableCell>
                <TableCell align="left">{item.seller}</TableCell>
                <TableCell align="left">{item.type}</TableCell>
                <TableCell align="left">
                  {format(new Date(item.date), 'dd/MM/yyyy HH:mm:ss')}
                </TableCell>
                <TableCell align="left">{item.batchID}</TableCell>
              </TableRow>
            )
          })}
          <Box sx={{ position: 'absolute', right: 0 }}>
            <Box sx={{ mt: 2 }}>
              <Divider textAlign="center">Summary</Divider>
            </Box>
            <Box>
              <TableRow>
                <TableCell>Credits</TableCell>
                <TableCell align="right">
                  {summary.credits.toLocaleString('pt-br', {
                    style: 'currency',
                    currency: 'BRL',
                  })}
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell>Debts</TableCell>
                <TableCell align="right">
                  {summary.debts.toLocaleString('pt-br', {
                    style: 'currency',
                    currency: 'BRL',
                  })}
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell>Total</TableCell>
                <TableCell align="right">
                  {summary.total.toLocaleString('pt-br', {
                    style: 'currency',
                    currency: 'BRL',
                  })}
                </TableCell>
              </TableRow>
            </Box>
          </Box>
        </TableBody>
      </Table>
    </TableContainer>
  )
}