import { Container, Image, Text, VStack } from "@chakra-ui/react";

export const Header = () => {
    return (
        <Container
            maxWidth={{ base: "100%" }}
            backgroundColor={"experiview.emerald"}
        >
            <Container as="header" width={"80%"}>
                <VStack p="1rem">
                    <Image src="/novologo.svg" />
                    <Text
                        color="experiview.blue"
                        fontWeight={"bold"}
                        fontSize={"2rem"}
                    >
                        Experiview
                    </Text>
                </VStack>
            </Container>
        </Container>
    );
};
