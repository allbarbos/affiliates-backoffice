import Button from '@mui/material/Button';
import Stack from '@mui/material/Stack';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import { Link } from "react-router-dom";

import Header from '../../components/Header'

export const Home = () => {
    return (
        <>
            <Header />
            <main>
                <Box sx={{ bgcolor: 'background.paper', pt: 8, pb: 6 }}>
                    <Container maxWidth="sm">
                        <Typography
                            component="h1"
                            variant="h2"
                            align="center"
                            color="text.primary"
                            gutterBottom
                        >
                            Afiliate Backoffice
                        </Typography>
                        <Typography variant="h5" align="center" color="text.secondary" paragraph>
                            Manage your sales!
                        </Typography>
                        <Stack
                            sx={{ pt: 4 }}
                            direction="row"
                            spacing={2}
                            justifyContent="center"
                        >
                            <Link to="upload" style={{ textDecoration: 'none' }}>
                                <Button variant="contained">File upload</Button>
                            </Link>
                            <Link to="transaction" style={{ textDecoration: 'none' }}>
                                <Button variant="outlined">Transactions</Button>
                            </Link>
                        </Stack>
                    </Container>
                </Box>
            </main>
        </>
    );
}