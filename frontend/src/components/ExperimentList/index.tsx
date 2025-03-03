import { FC } from "react";
import { Experiment } from "../../types";
import { Box, Flex, Heading, Text, VStack } from "@chakra-ui/react";

interface ExperimentListProps {
    experiments: Experiment[];
}
export const ExperimentList: FC<ExperimentListProps> = ({ experiments }) => {
    return (
        <VStack spacing={4} align="stretch" mt={4}>
            {experiments.map((experiment) => (
                <Box
                    key={experiment.id}
                    border="1px"
                    borderColor={"experiview.emerald"}
                    borderRadius="md"
                    p={4}
                    _hover={{ bg: "experiview.blue" }}
                    transition="background-color 0.2s"
                >
                    <Flex justify="space-between" align="start">
                        <Heading as="h3" size="md" color={"black"}>
                            {experiment.name}
                        </Heading>
                        <Text fontSize="sm" color="black">
                            {new Date(experiment.date).toLocaleDateString()}
                        </Text>
                    </Flex>
                    <Text mt={2} color={"black"}>
                        {experiment.description}
                    </Text>
                </Box>
            ))}
        </VStack>
    );
};
