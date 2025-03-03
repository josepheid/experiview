import "@testing-library/jest-dom";

import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { Experiments } from "../Experiments";
import { describe, it, expect, beforeEach } from "vitest";

import { vi } from "vitest";

// Mock fetch API
vi.stubGlobal("fetch", vi.fn());

describe("Experiments Component", () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("renders component correctly", async () => {
        render(<Experiments />);
        expect(screen.getByText(/Experiments/i)).toBeInTheDocument();
    });

    it("fetches and displays experiments", async () => {
        const mockExperiments = [
            {
                id: "1",
                name: "Test Experiment 1",
                date: "2024-03-01",
                description: "Desc 1",
            },
            {
                id: "2",
                name: "Test Experiment 2",
                date: "2024-02-15",
                description: "Desc 2",
            },
        ];

        vi.mocked(fetch).mockResolvedValueOnce(
            new Response(JSON.stringify(mockExperiments), {
                status: 200,
                headers: { "Content-Type": "application/json" },
            })
        );

        render(<Experiments />);
        await waitFor(() => expect(fetch).toHaveBeenCalled());

        expect(screen.getByText("Test Experiment 1")).toBeInTheDocument();
        expect(screen.getByText("Test Experiment 2")).toBeInTheDocument();
    });

    it("handles sorting experiments", async () => {
        const mockExperiments = [
            {
                id: "1",
                name: "B Experiment",
                date: "2024-03-01",
                description: "Desc 1",
            },
            {
                id: "2",
                name: "A Experiment",
                date: "2024-02-15",
                description: "Desc 2",
            },
        ];

        vi.mocked(fetch).mockResolvedValueOnce(
            new Response(JSON.stringify(mockExperiments), {
                status: 200,
                headers: { "Content-Type": "application/json" },
            })
        );

        render(<Experiments />);
        await waitFor(() => expect(fetch).toHaveBeenCalled());

        const sortSelect = screen.getByRole("combobox");
        fireEvent.change(sortSelect, { target: { value: "titleAsc" } });

        await waitFor(() => {
            const firstExperiment = screen.getAllByRole("heading")[0];
            expect(firstExperiment).toHaveTextContent("A Experiment");
        });
    });
});
