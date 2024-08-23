import {
  createContext,
  useContext,
  useState,
  FC,
  ReactNode,
  ReactElement,
} from "react";

/**
 * Represents the structure of the SearchJobsContext.
 *
 * @property {string} searchTerm - The current search term used to filter job vacancies.
 * @property {function} setSearchTerm - Function to update the search term.
 */
interface SearchJobsContextType {
  searchTerm: string;
  setSearchTerm: (term: string) => void;
}

interface SearchJobsProviderProps {
  children: ReactNode;
}

const SearchJobsContext = createContext<SearchJobsContextType | undefined>(
  undefined,
);

/**
 * SearchJobsProvider component that manages and provides the search term state.
 * It allows components to access and update the search term used for filtering job vacancies.
 *
 * @param {SearchJobsProviderProps} props - The props for the SearchJobsProvider component.
 * @returns {ReactElement} The rendered SearchJobsProvider component.
 */
export const SearchJobsProvider: FC<SearchJobsProviderProps> = ({
  children,
}: SearchJobsProviderProps): ReactElement => {
  const [searchTerm, setSearchTerm] = useState<string>("");

  return (
    <SearchJobsContext.Provider value={{ searchTerm, setSearchTerm }}>
      {children}
    </SearchJobsContext.Provider>
  );
};

/**
 * Custom hook to access the SearchJobsContext.
 * Throws an error if used outside of a SearchJobsProvider.
 *
 * @returns {SearchJobsContextType} The context value containing searchTerm and setSearchTerm function.
 */
export const useSearchJobs = (): SearchJobsContextType => {
  const context: SearchJobsContextType | undefined =
    useContext(SearchJobsContext);
  if (!context) {
    throw new Error("useSearchJobs must be used within a SearchJobsProvider");
  }
  return context;
};
