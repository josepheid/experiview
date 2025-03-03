import { Container } from "@chakra-ui/react";
import { Header } from "./components/Header";
import { Experiments } from "./components/Experiments";

function App() {
    return (
        <main>
            <Container maxWidth={{ base: "100%" }} minH={"100vh"} px="0">
                <Header />
                <Experiments />
            </Container>
        </main>
    );
}

export default App;
