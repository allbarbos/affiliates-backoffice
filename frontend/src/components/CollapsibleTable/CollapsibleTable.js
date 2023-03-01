import {
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
} from '@mui/material'

import { Row } from './Row'

export const CollapsibleTable = ({ items }) => {
  return (
    <>
      {items?.length ? (
        <>
          <TableContainer component={Paper}>
            <Table aria-label="collapsible table">
              <TableHead>
                <TableRow>
                  <TableCell />
                  <TableCell align="left">ID</TableCell>
                  <TableCell align="left">Affiliate ID</TableCell>
                  <TableCell align="left">Status</TableCell>
                  <TableCell align="left">Created At</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {items.map((item) => (
                  <Row key={item.id} row={item} />
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </>
      ) : (
        <></>
      )}
    </>
  )
}
