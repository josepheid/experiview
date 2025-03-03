import {
    Alert,
    AlertIcon,
    Button,
    FormControl,
    FormLabel,
    Input,
    Modal,
    ModalBody,
    ModalCloseButton,
    ModalContent,
    ModalFooter,
    ModalHeader,
    ModalOverlay,
    Textarea,
    VStack,
} from "@chakra-ui/react";
import { FC, useState } from "react";
import { Experiment } from "../../types";

interface AddExperimentProps {
    isOpen: boolean;
    onClose: () => void;
    onAddExperiment: (experiment: Experiment) => void;
}

export const AddExperiment: FC<AddExperimentProps> = ({
    isOpen,
    onClose,
    onAddExperiment,
}) => {
    const [name, setName] = useState("");
    const [date, setDate] = useState(new Date().toISOString().split("T")[0]);
    const [description, setDescription] = useState("");
    const [error, setError] = useState("");

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();

        if (!name.trim() || !date || !description.trim()) {
            setError("All fields are required");
            return;
        }

        const newExperiment: Experiment = {
            id: "", // Will be set by the parent component
            name: name.trim(),
            date,
            description: description.trim(),
        };

        onAddExperiment(newExperiment);

        // Reset form
        setName("");
        setDescription("");
        setDate(new Date().toISOString().split("T")[0]);
        setError("");
    };
    return (
        <Modal isOpen={isOpen} onClose={onClose}>
            <ModalOverlay />
            <ModalContent>
                <ModalHeader>Add Experiment</ModalHeader>
                <ModalCloseButton />
                <ModalBody>
                    <VStack
                        as="form"
                        spacing={4}
                        align="stretch"
                        onSubmit={handleSubmit}
                    >
                        {error && (
                            <Alert status="error" borderRadius="md">
                                <AlertIcon />
                                {error}
                            </Alert>
                        )}

                        <FormControl isRequired>
                            <FormLabel htmlFor="name">
                                Experiment Name
                            </FormLabel>
                            <Input
                                id="name"
                                data-testid="name"
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                                focusBorderColor="experiview.emerald"
                            />
                        </FormControl>

                        <FormControl isRequired>
                            <FormLabel htmlFor="date">Date</FormLabel>
                            <Input
                                type="date"
                                id="date"
                                data-testid="date"
                                value={date}
                                onChange={(e) => setDate(e.target.value)}
                                focusBorderColor="experiview.emerald"
                            />
                        </FormControl>

                        <FormControl isRequired>
                            <FormLabel htmlFor="description">
                                Description
                            </FormLabel>
                            <Textarea
                                id="description"
                                data-testid="description"
                                value={description}
                                onChange={(e) => setDescription(e.target.value)}
                                rows={4}
                                focusBorderColor="experiview.emerald"
                            />
                        </FormControl>

                        <Button
                            type="submit"
                            w="full"
                            colorScheme="green"
                            bg={"experiview.blue"}
                            _hover={{ bg: "experiview.lightblue" }}
                        >
                            Register Experiment
                        </Button>
                    </VStack>
                </ModalBody>

                <ModalFooter>
                    <Button
                        color={"white"}
                        bg={"experiview.blue"}
                        _hover={{ bg: "experiview.lightblue" }}
                        mr={3}
                        onClick={onClose}
                    >
                        Close
                    </Button>
                </ModalFooter>
            </ModalContent>
        </Modal>
    );
};
