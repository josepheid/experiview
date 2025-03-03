import { render, screen, fireEvent } from "@testing-library/react";
import { ChakraProvider } from "@chakra-ui/react";
import { describe, it, expect } from "vitest";
import { ExperimentList } from ".";

const mockExperiments = [
    {
        id: "1",
        name: "Experiment 1",
        date: "2024-02-28",
        description: "First experiment description",
    },
    {
        id: "2",
        name: "Experiment 2",
        date: "2024-02-27",
        description: "Second experiment description",
    },
    {
        id: "3",
        name: "Experiment 3",
        date: "2024-02-26",
        description: "Third experiment description",
    },
    {
        id: "4",
        name: "Experiment 4",
        date: "2024-02-25",
        description: "Fourth experiment description",
    },
    {
        id: "5",
        name: "Experiment 5",
        date: "2024-02-24",
        description: "Fifth experiment description",
    },
    {
        id: "6",
        name: "Experiment 6",
        date: "2024-02-23",
        description: "Sixth experiment description",
    },
];

describe("ExperimentList Component", () => {
    const renderComponent = () =>
        render(
            <ChakraProvider>
                <ExperimentList experiments={mockExperiments} />
            </ChakraProvider>
        );

    it("displays the first page of experiments", () => {
        renderComponent();
        expect(screen.getByText("Experiment 1")).toBeInTheDocument();
        expect(screen.getByText("Experiment 5")).toBeInTheDocument();
        expect(screen.queryByText("Experiment 6")).not.toBeInTheDocument();
    });

    it("pagination works correctly", () => {
        renderComponent();
        const nextButton = screen.getByRole("button", { name: /next/i });
        fireEvent.click(nextButton);
        expect(screen.getByText("Experiment 6")).toBeInTheDocument();
        expect(screen.queryByText("Experiment 1")).not.toBeInTheDocument();
    });

    it("clicking an experiment opens modal with details", () => {
        renderComponent();
        fireEvent.click(screen.getByText("Experiment 1"));
        expect(
            screen.getByTestId(`experiment-description`)
        ).toBeInTheDocument();
        expect(screen.getByTestId("experiment-id")).toBeInTheDocument();
    });

    it("closing modal resets selection", () => {
        renderComponent();
        fireEvent.click(screen.getByText("Experiment 1"));
        fireEvent.click(screen.getByText("Close"));
        expect(screen.queryByText("ID: 1")).not.toBeInTheDocument();
    });
});
