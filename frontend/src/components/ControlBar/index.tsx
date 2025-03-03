import { AddIcon, SearchIcon } from "@chakra-ui/icons";
import {
    Box,
    Button,
    Flex,
    FormControl,
    FormLabel,
    Input,
    InputGroup,
    InputLeftElement,
    Select,
    Stack,
    useDisclosure,
} from "@chakra-ui/react";
import { AddExperiment } from "../AddExperiment";
import { FC, useEffect, useState } from "react";
import { Experiment } from "../../types";

interface ControlBarProps {
    filterText?: string;
    onFilterChange: (text: string) => void;
    sortOption?: string;
    onSortChange: (option: string) => void;
    startDate?: string;
    endDate?: string;
    onStartDateChange: (date: string) => void;
    onEndDateChange: (date: string) => void;
}

export const ControlBar: FC<ControlBarProps> = ({
    filterText,
    onFilterChange,
    sortOption,
    onSortChange,
    startDate,
    endDate,
    onStartDateChange,
    onEndDateChange,
}) => {
    const { isOpen, onOpen, onClose } = useDisclosure();
    const [defaultEndDate, setDefaultEndDate] = useState("");

    useEffect(() => {
        const today = new Date().toISOString().split("T")[0];
        setDefaultEndDate(today);
        if (!endDate) {
            onEndDateChange(today);
        }
    }, [endDate, onEndDateChange]);

    const apiURL = `${import.meta.env.VITE_API_URL}/experiview/experiments`;
    const apiKey = import.meta.env.VITE_API_KEY;
    console.log(apiKey);

    const addExperiment = async (experiment: Experiment) => {
        const response = await fetch(apiURL as string, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "x-api-key": apiKey as string,
            },
            body: JSON.stringify(experiment),
        });
        console.log(response);
    };

    return (
        <>
            <Stack
                direction={{ base: "column", sm: "row" }}
                spacing={4}
                align="center"
                justify="space-between"
                mb={6}
            >
                <InputGroup maxW={{ base: "full", sm: "30%" }}>
                    <InputLeftElement pointerEvents="none">
                        <SearchIcon color="gray.400" />
                    </InputLeftElement>
                    <Input
                        type="text"
                        placeholder="Filter experiments..."
                        value={filterText}
                        onChange={(e) => onFilterChange(e.target.value)}
                        focusBorderColor="experiview.emerald"
                    />
                </InputGroup>
                <Flex
                    direction={{ base: "column", sm: "row" }}
                    gap={3}
                    align={{ base: "stretch", sm: "center" }}
                    w={{ base: "full", md: "auto" }}
                >
                    <Box>
                        <FormControl size="sm">
                            <FormLabel fontSize="xs" mb={1}>
                                From
                            </FormLabel>
                            <Input
                                type="date"
                                value={startDate}
                                onChange={(e) =>
                                    onStartDateChange(e.target.value)
                                }
                                focusBorderColor="experiview.emerald"
                                size="sm"
                                max={endDate || defaultEndDate}
                            />
                        </FormControl>
                    </Box>
                    <Box>
                        <FormControl size="sm">
                            <FormLabel fontSize="xs" mb={1}>
                                To
                            </FormLabel>
                            <Input
                                type="date"
                                value={endDate || defaultEndDate}
                                onChange={(e) =>
                                    onEndDateChange(e.target.value)
                                }
                                focusBorderColor="experiview.emerald"
                                size="sm"
                                max={defaultEndDate}
                            />
                        </FormControl>
                    </Box>
                </Flex>
                <Select
                    maxW={{ base: "full", sm: "30%" }}
                    value={sortOption}
                    onChange={(e) => onSortChange(e.target.value)}
                    focusBorderColor="experiview.emerald"
                >
                    <option value="dateDesc">Newest First</option>
                    <option value="dateAsc">Oldest First</option>
                    <option value="titleAsc">Title (A-Z)</option>
                    <option value="titleDesc">Title (Z-A)</option>
                </Select>
            </Stack>
            <Button
                backgroundColor={"experiview.blue"}
                color="white"
                _hover={{
                    backgroundColor: "experiview.lightblue",
                    color: "white",
                }}
                onClick={onOpen}
            >
                <AddIcon mr={"0.5rem"} />
                New Experiment
            </Button>
            <AddExperiment
                isOpen={isOpen}
                onClose={onClose}
                onAddExperiment={addExperiment}
            />
        </>
    );
};
