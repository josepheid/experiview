import { AddIcon, SearchIcon } from "@chakra-ui/icons";
import {
    Button,
    Input,
    InputGroup,
    InputLeftElement,
    Select,
    Stack,
    useDisclosure,
} from "@chakra-ui/react";
import { AddExperiment } from "../AddExperiment";

export const ControlBar = () => {
    const { isOpen, onOpen, onClose } = useDisclosure();

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
                        // value={filterText}
                        // onChange={(e) => onFilterChange(e.target.value)}
                        focusBorderColor="experiview.emerald"
                    />
                </InputGroup>

                <Select
                    maxW={{ base: "full", sm: "30%" }}
                    // value={sortOption}
                    // onChange={(e) => onSortChange(e.target.value)}
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
            <AddExperiment isOpen={isOpen} onClose={onClose} />
        </>
    );
};
