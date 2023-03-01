import { useState, useEffect } from 'react';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import { ReactNotifications } from 'react-notifications-component'
import 'react-notifications-component/dist/theme.css'

import Header from '../../components/Header'
import Toast from '../../components/Toast'
import SpanningTable from '../../components/SpanningTable'
import * as AffiliateAPI from '../../services/affiliate-api'

export const Transaction = () => {
    const [items, setItems] = useState([]);
    const [summary, setSummary] = useState({
        debts: 0,
        credits: 0,
        total: 0
    });

    async function fetchData() {
        try {
            const { data } = await AffiliateAPI.GetTransactions();
            setItems(data)
            setSummary(data.reduce((acc, { value }) => {
                value >= 0 ? acc.credits += value : acc.debts += value
                acc.total += value
                return acc
            }, { credits: 0, debts: 0, total: 0 }))
        } catch (error) {
            return Toast({
                type: "danger",
                title: "Oops!",
                message: error.message,
            })
        }
    }

    useEffect(() => {
        fetchData();
    }, []);

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
                <Typography variant="h5" align="center" color="text.secondary" paragraph>
                    Imported sales
                </Typography>
                <Box
                    sx={{
                        marginTop: 8,
                        display: 'flex',
                        flexDirection: 'column',
                        alignItems: 'center',
                        justifyContent: 'center',
                        width: "75%"
                    }}
                >
                    {
                        items.length ?
                            <SpanningTable items={items} summary={summary} /> :
                            <Typography variant="h7" align="center" color="text.secondary" paragraph>
                                No imported files
                            </Typography>
                    }
                </Box>
            </Box>

        </>
    )
}