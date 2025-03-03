import { Container, Text } from "@chakra-ui/react";
import { ControlBar } from "../ControlBar";
import { useEffect, useState } from "react";
import { Experiment } from "../../types";
import { ExperimentList } from "../ExperimentList";

export const Experiments = () => {
    const [experiments, setExperiments] = useState<Experiment[]>([
        { date: "2024/04/03", description: "yeah", id: "1013", name: "drug" },
        { date: "2024/04/06", description: "no", id: "1014", name: "drug" },
        { date: "2024/04/08", description: "john", id: "1015", name: "drug" },
        { date: "2024/04/12", description: "deer", id: "1016", name: "drug" },
        { date: "2024/04/09", description: "cash", id: "1017", name: "drug" },
        { date: "2024/04/08", description: "john", id: "1018", name: "amber" },
        { date: "2024/04/12", description: "deer", id: "1019", name: "green" },
        { date: "2024/04/09", description: "cash", id: "1020", name: "red" },
    ]);
    const [startDate, setStartDate] = useState<string | undefined>(undefined);
    const [endDate, setEndDate] = useState<string | undefined>(undefined);
    const [filter, setFilter] = useState<string | undefined>(undefined);
    const [sort, setSort] = useState<string | undefined>(undefined);
    const [shouldReload, setShouldReload] = useState(false);

    const fetchExperiments = async (
        filter: string | undefined,
        startDate: string | undefined,
        endDate: string | undefined
    ) => {
        const apiKey = import.meta.env.VITE_API_KEY;
        let apiURL = `${import.meta.env.VITE_API_URL}/experiview/experiments`;

        const queryParams = new URLSearchParams();
        if (filter) queryParams.append("filter", filter);
        if (startDate) queryParams.append("startDate", startDate);
        if (endDate) queryParams.append("endDate", endDate);

        if (queryParams.toString()) {
            apiURL += `?${queryParams.toString()}`;
        }

        try {
            const resp = await fetch(apiURL, {
                headers: {
                    "Content-Type": "application/json",
                    "x-api-key": apiKey as string,
                },
            });

            const data = await resp.json();

            setExperiments(data);
        } catch (error) {
            console.error("Error fetching experiments:", error);
        }

        setShouldReload(false);
    };

    useEffect(() => {
        if (filter || startDate || endDate)
            fetchExperiments(filter, startDate, endDate);
    }, [filter, startDate, endDate, shouldReload]);

    const onStartDateChange = (val: string) => {
        setStartDate(val);
    };

    const onEndDateChange = (val: string) => {
        setEndDate(val);
    };

    const onFilterChange = (val: string) => {
        setFilter(val);
    };

    const onSortChange = (val: string) => {
        setSort(val);
        const updatedExperiments = [...experiments];
        switch (val) {
            case "titleAsc":
                updatedExperiments.sort((a, b) => a.name.localeCompare(b.name));
                break;
            case "titleDesc":
                updatedExperiments.sort((a, b) => b.name.localeCompare(a.name));
                break;
            case "dateAsc":
                updatedExperiments.sort(
                    (a, b) =>
                        new Date(a.date).getTime() - new Date(b.date).getTime()
                );
                break;
            case "dateDesc":
                updatedExperiments.sort(
                    (a, b) =>
                        new Date(b.date).getTime() - new Date(a.date).getTime()
                );
                break;
        }

        setExperiments(updatedExperiments);
    };

    return (
        <Container py="1rem" minW={{ base: "100%", md: "80%" }}>
            <Text fontSize={"2rem"} fontWeight={"bold"} py="1rem">
                Experiments
            </Text>
            <ControlBar
                startDate={startDate}
                onStartDateChange={onStartDateChange}
                endDate={endDate}
                onEndDateChange={onEndDateChange}
                filterText={filter}
                onFilterChange={onFilterChange}
                sortOption={sort}
                onSortChange={onSortChange}
                setShouldReload={setShouldReload}
            />
            <ExperimentList experiments={experiments} />
        </Container>
    );
};
