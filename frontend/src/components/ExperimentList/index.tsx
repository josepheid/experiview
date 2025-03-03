import { FC, useState } from "react";
import { Experiment } from "../../types";
import {
    Box,
    Button,
    Flex,
    Heading,
    Modal,
    ModalBody,
    ModalCloseButton,
    ModalContent,
    ModalFooter,
    ModalHeader,
    ModalOverlay,
    Text,
    useDisclosure,
    VStack,
} from "@chakra-ui/react";

interface ExperimentListProps {
    experiments: Experiment[];
}

const ITEMS_PER_PAGE = 5;

export const ExperimentList: FC<ExperimentListProps> = ({ experiments }) => {
    const { isOpen, onOpen, onClose } = useDisclosure();
    const [selected, setSelected] = useState<Experiment | undefined>(undefined);
    const [currentPage, setCurrentPage] = useState(1);

    const totalPages = Math.ceil(experiments.length / ITEMS_PER_PAGE);
    const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
    const currentExperiments = experiments.slice(
        startIndex,
        startIndex + ITEMS_PER_PAGE
    );

    return (
        <>
            <VStack spacing={4} align="stretch" mt={4}>
                {currentExperiments.map((experiment) => (
                    <Box
                        as={"button"}
                        key={experiment.id}
                        border="1px"
                        borderColor={"experiview.emerald"}
                        borderRadius="md"
                        p={4}
                        _hover={{ bg: "experiview.emerald" }}
                        transition="background-color 0.2s"
                        onClick={() => {
                            setSelected(experiment);
                            onOpen();
                        }}
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

            <Flex justify="center" mt={4}>
                <Button
                    color={"white"}
                    onClick={() =>
                        setCurrentPage((prev) => Math.max(prev - 1, 1))
                    }
                    isDisabled={currentPage === 1}
                    bg={"experiview.blue"}
                    _hover={{ bg: "experiview.lightblue" }}
                    minW={"6.25rem"}
                >
                    Previous
                </Button>
                <Text fontSize="lg" mx={4}>
                    Page {currentPage} of {totalPages}
                </Text>
                <Button
                    color={"white"}
                    onClick={() =>
                        setCurrentPage((prev) => Math.min(prev + 1, totalPages))
                    }
                    isDisabled={currentPage === totalPages}
                    bg={"experiview.blue"}
                    _hover={{ bg: "experiview.lightblue" }}
                    minW={"6.25rem"}
                >
                    Next
                </Button>
            </Flex>

            {selected && (
                <Modal isOpen={isOpen} onClose={onClose}>
                    <ModalOverlay />
                    <ModalContent>
                        <ModalHeader>{selected.name}</ModalHeader>
                        <ModalCloseButton />
                        <ModalBody>
                            <VStack align="start">
                                <Text>
                                    <strong>ID:</strong> {selected.id}
                                </Text>
                                <Text>
                                    <strong>Date:</strong>{" "}
                                    {new Date(
                                        selected.date
                                    ).toLocaleDateString()}
                                </Text>
                                <Text>
                                    <strong>Description:</strong>{" "}
                                    {selected.description}
                                </Text>
                            </VStack>
                        </ModalBody>

                        <ModalFooter>
                            <Button
                                color={"white"}
                                bg={"experiview.blue"}
                                _hover={{ bg: "experiview.lightblue" }}
                                onClick={() => {
                                    onClose();
                                    setSelected(undefined);
                                }}
                            >
                                Close
                            </Button>
                        </ModalFooter>
                    </ModalContent>
                </Modal>
            )}
        </>
    );
};
