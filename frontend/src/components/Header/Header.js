import { AppBar, Toolbar, Typography } from '@mui/material';
import { Link } from "react-router-dom";

export const Header = () => {
    return (
        <AppBar position="relative">
            <Toolbar>
                <Link to="/" style={{ textDecoration: 'none' }}>
                    <Typography variant="h6" color="white" noWrap>
                        Afiliate Backoffice
                    </Typography>
                </Link>
            </Toolbar>
        </AppBar>
    )
}