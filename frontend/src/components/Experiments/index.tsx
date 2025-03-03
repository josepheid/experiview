import { Container, Text } from "@chakra-ui/react";
import { ControlBar } from "../ControlBar";
import { useState } from "react";
import { Experiment } from "../../types";
import { ExperimentList } from "../ExperimentList";

export const Experiments = () => {
    const [experiments, setExperiments] = useState<Experiment[]>([
        { date: "2024/04/03", description: "yeah", id: "1013", name: "drug" },
    ]);

    return (
        <Container py="1rem" minW={{ base: "100%", md: "80%" }}>
            <Text fontSize={"2rem"} fontWeight={"bold"} py="1rem">
                Experiments
            </Text>
            <ControlBar />
            <ExperimentList experiments={experiments} />
        </Container>
    );
};
