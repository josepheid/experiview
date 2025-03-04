import { Container, Text, useToast } from "@chakra-ui/react";
import { ControlBar } from "../ControlBar";
import { useEffect, useState } from "react";
import { Experiment } from "../../types";
import { ExperimentList } from "../ExperimentList";

export const Experiments = () => {
    const [experiments, setExperiments] = useState<Experiment[]>([]);
    const [startDate, setStartDate] = useState<string | undefined>(undefined);
    const [endDate, setEndDate] = useState<string | undefined>(undefined);
    const [filter, setFilter] = useState<string | undefined>(undefined);
    const [debouncedFilter, setDebouncedFilter] = useState<string | undefined>(
        undefined
    );
    const [sort, setSort] = useState<string>("dateDesc");
    const [shouldReload, setShouldReload] = useState(false);
    const toast = useToast();

    const sortExperiments = (experiments: Experiment[], sortOption: string) => {
        return [...experiments].sort((a, b) => {
            if (sortOption === "dateAsc") {
                return new Date(a.date).getTime() - new Date(b.date).getTime();
            } else if (sortOption === "dateDesc") {
                return new Date(b.date).getTime() - new Date(a.date).getTime();
            } else if (sortOption === "titleAsc") {
                return a.name.localeCompare(b.name);
            } else if (sortOption === "titleDesc") {
                return b.name.localeCompare(a.name);
            }
            return 0;
        });
    };

    useEffect(() => {
        const fetchExperiments = async (
            debouncedFilter: string | undefined,
            startDate: string | undefined,
            endDate: string | undefined
        ) => {
            const apiKey = import.meta.env.VITE_API_KEY;
            let apiURL = `${
                import.meta.env.VITE_API_URL
            }/experiview/experiments`;

            const queryParams = new URLSearchParams();
            if (debouncedFilter) queryParams.append("filter", debouncedFilter);
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

                let data: Experiment[] = await resp.json();

                data = sortExperiments(data, sort);

                setExperiments(data);
            } catch (error) {
                toast({
                    title: "Error",
                    description: `Something went wrong. Please try again. Error: ${error}`,
                    status: "error",
                    duration: 5000,
                    isClosable: true,
                    position: "top-right",
                });
            }

            setShouldReload(false);
        };
        if (shouldReload) {
            fetchExperiments(debouncedFilter, startDate, endDate).then(() =>
                setShouldReload(false)
            );
        }
    }, [debouncedFilter, startDate, endDate, shouldReload, sort, toast]);

    useEffect(() => {
        setExperiments((prev) => sortExperiments(prev, sort));
    }, [sort]);

    useEffect(() => {
        const handler = setTimeout(() => {
            setDebouncedFilter(filter); // Update after 2 seconds
            setShouldReload(true);
        }, 750);

        return () => clearTimeout(handler); // Clear timeout if user types again
    }, [filter]);

    const onStartDateChange = (val: string) => {
        setShouldReload(true);
        setStartDate(val);
    };

    const onEndDateChange = (val: string) => {
        setShouldReload(true);
        setEndDate(val);
    };

    const onFilterChange = (val: string) => {
        setFilter(val);
    };

    const onSortChange = (val: string) => {
        setSort(val);
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
