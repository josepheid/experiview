import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { ChakraProvider } from "@chakra-ui/react";
import { AddExperiment } from ".";
import { Experiment } from "../../types";

describe("AddExperiment Component", () => {
    let onCloseMock: () => void;
    let onAddExperimentMock: (experiment: Experiment) => void;

    beforeEach(() => {
        onCloseMock = vi.fn();
        onAddExperimentMock = vi.fn();
    });

    const setup = () =>
        render(
            <ChakraProvider>
                <AddExperiment
                    isOpen={true}
                    onClose={onCloseMock}
                    onAddExperiment={onAddExperimentMock}
                />
            </ChakraProvider>
        );

    it("renders correctly", () => {
        setup();
        expect(screen.getByText("Add Experiment")).toBeInTheDocument();
        expect(screen.getByText("Experiment Name")).toBeInTheDocument();
        expect(screen.getByText("Date")).toBeInTheDocument();
        expect(screen.getByText("Description")).toBeInTheDocument();
        expect(screen.getByText("Register Experiment")).toBeInTheDocument();
    });

    it("allows input fields to be updated", () => {
        setup();

        const nameInput = screen.getByTestId("name");
        fireEvent.change(nameInput, { target: { value: "Test Experiment" } });
        expect(nameInput).toHaveValue("Test Experiment");

        const descriptionInput = screen.getByTestId("description");
        fireEvent.change(descriptionInput, {
            target: { value: "This is a test" },
        });
        expect(descriptionInput).toHaveValue("This is a test");
    });

    it("calls onAddExperiment with correct data on valid submit", () => {
        setup();

        fireEvent.change(screen.getByTestId("name"), {
            target: { value: "Test Exp" },
        });
        fireEvent.change(screen.getByTestId("date"), {
            target: { value: "2024-03-01" },
        });
        fireEvent.change(screen.getByTestId("description"), {
            target: { value: "A description" },
        });

        const submitButton = screen.getByText("Register Experiment");
        fireEvent.click(submitButton);

        expect(onAddExperimentMock).toHaveBeenCalledTimes(1);
        expect(onAddExperimentMock).toHaveBeenCalledWith({
            id: "",
            name: "Test Exp",
            date: "2024-03-01",
            description: "A description",
        });
    });

    it("closes the modal when Close button is clicked", () => {
        setup();

        const closeButton = screen.getByText("Close");
        fireEvent.click(closeButton);

        expect(onCloseMock).toHaveBeenCalledTimes(1);
    });
});
