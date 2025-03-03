import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { ControlBar } from "../ControlBar";
import { ChakraProvider } from "@chakra-ui/react";

// Mocking toast
vi.mock("@chakra-ui/react", async () => {
    const actual = await vi.importActual("@chakra-ui/react");
    return {
        ...actual,
        useToast: () => () => {},
    };
});

describe("ControlBar", () => {
    let onFilterChange: () => void,
        onSortChange: () => void,
        onStartDateChange: () => void,
        onEndDateChange: () => void,
        setShouldReload: () => void;

    beforeEach(() => {
        onFilterChange = vi.fn();
        onSortChange = vi.fn();
        onStartDateChange = vi.fn();
        onEndDateChange = vi.fn();
        setShouldReload = vi.fn();
    });

    const renderComponent = () =>
        render(
            <ChakraProvider>
                <ControlBar
                    filterText=""
                    onFilterChange={onFilterChange}
                    sortOption="dateDesc"
                    onSortChange={onSortChange}
                    startDate=""
                    endDate=""
                    onStartDateChange={onStartDateChange}
                    onEndDateChange={onEndDateChange}
                    setShouldReload={setShouldReload}
                />
            </ChakraProvider>
        );

    it("renders input fields and button", () => {
        renderComponent();
        expect(
            screen.getByPlaceholderText("Filter experiments...")
        ).toBeInTheDocument();
        expect(screen.getByText("New Experiment")).toBeInTheDocument();
    });

    it("calls onFilterChange when filter input is changed", () => {
        renderComponent();
        const filterInput = screen.getByPlaceholderText(
            "Filter experiments..."
        );
        fireEvent.change(filterInput, { target: { value: "test" } });
        expect(onFilterChange).toHaveBeenCalledWith("test");
    });

    it("calls onSortChange when sort option is selected", () => {
        renderComponent();
        const select = screen.getByRole("combobox");
        fireEvent.change(select, { target: { value: "titleAsc" } });
        expect(onSortChange).toHaveBeenCalledWith("titleAsc");
    });

    it("calls onStartDateChange when start date is changed", () => {
        renderComponent();
        const startDateInput = screen.getByLabelText("From");
        fireEvent.change(startDateInput, { target: { value: "2024-03-01" } });
        expect(onStartDateChange).toHaveBeenCalledWith("2024-03-01");
    });

    it("calls onEndDateChange when end date is changed", () => {
        renderComponent();
        const endDateInput = screen.getByLabelText("To");
        fireEvent.change(endDateInput, { target: { value: "2024-03-05" } });
        expect(onEndDateChange).toHaveBeenCalledWith("2024-03-05");
    });

    it("opens the AddExperiment modal when 'New Experiment' button is clicked", async () => {
        renderComponent();
        fireEvent.click(screen.getByText("New Experiment"));
        await waitFor(() =>
            expect(screen.getByText("Add Experiment")).toBeInTheDocument()
        );
    });
});
