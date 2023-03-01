import { useState } from 'react';
import Box from '@mui/material/Box';
import Collapse from '@mui/material/Collapse';
import IconButton from '@mui/material/IconButton';
import TableCell from '@mui/material/TableCell';
import TableRow from '@mui/material/TableRow';
import Typography from '@mui/material/Typography';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableHead from '@mui/material/TableHead';
import { format } from 'date-fns/esm'

export const Row = ({ row }) => {
    const [open, setOpen] = useState(false);

    return (
        <>
            <TableRow sx={{ '& > *': { borderBottom: 'unset' } }}>
                <TableCell>
                    <IconButton
                        aria-label="expand row"
                        size="small"
                        onClick={() => setOpen(!open)}
                    >
                        {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
                    </IconButton>
                </TableCell>
                <TableCell component="th" scope="row">
                    {row.id}
                </TableCell>
                <TableCell align="left">{row.affiliateID}</TableCell>
                <TableCell align="left">{row.status}</TableCell>
                <TableCell align="left">{format(new Date(row.createdAt), 'dd/MM/yyyy HH:mm:ss')}</TableCell>
            </TableRow>
            <TableRow>
                <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        <Box sx={{ margin: 1 }}>
                            {
                                row.errors?.length ?
                                    <>
                                        <Typography variant="h6" gutterBottom component="div">
                                            Errors
                                        </Typography>
                                        <Table size="small" aria-label="purchases">
                                            <TableHead>
                                                <TableRow>
                                                    <TableCell>Row</TableCell>
                                                    <TableCell>Errors</TableCell>
                                                </TableRow>
                                            </TableHead>
                                            <TableBody>
                                                {row.errors.sort((x, y) => x.row - y.row).map((error) => (
                                                    <TableRow key={error.row}>
                                                        <TableCell component="th" scope="row">
                                                            {error.row + 1}
                                                        </TableCell>
                                                        <TableCell>
                                                            {error.errors.map((err) => {
                                                                return (
                                                                    <span key={err}>{err}<br></br></span>
                                                                )
                                                            })}
                                                        </TableCell>

                                                    </TableRow>
                                                ))}
                                            </TableBody>
                                        </Table>
                                    </> :
                                    <Typography variant="h6" gutterBottom component="div">
                                        No errors
                                    </Typography>
                            }
                        </Box>
                    </Collapse>
                </TableCell>
            </TableRow>
        </>
    );
}