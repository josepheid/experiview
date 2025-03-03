import { Box, Container, HStack, Image, Text } from "@chakra-ui/react";

export const Header = () => {
    return (
        <Container
            maxWidth={{ base: "100%" }}
            backgroundColor={"experiview.emerald"}
        >
            <Box as="header">
                <HStack p="1rem">
                    <Image src="/novologo.svg" />
                    <Text
                        color="experiview.blue"
                        fontWeight={"bold"}
                        fontSize={"2rem"}
                    >
                        Experiview
                    </Text>
                </HStack>
            </Box>
        </Container>
    );
};
