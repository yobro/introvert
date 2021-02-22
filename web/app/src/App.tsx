import React from 'react';
import { Container, Navbar } from '../node_modules/react-bootstrap/esm/index';
import 'bootstrap/scss/bootstrap.scss';

function App() {
    return (
        <Container>
            <Navbar bg="dark" variant="dark">
               <Navbar.Brand>Introvert</Navbar.Brand> 
            </Navbar>
        </Container>
    );
}

export default App;