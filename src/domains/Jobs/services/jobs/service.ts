import axiosInstance from "../axios/axios";
import JwtService from "../jwt/service";
import type { JobsResponse } from "./response";

/**
 * A singleton service to manage job related API calls.
 */
class JobService {
  private static instance: JobService;

  /**
   * Private constructor to enforce the singleton pattern.
   */
  private constructor() {}

  /**
   * Returns the single instance of the JobService class.
   * If no instance exists, it creates one.
   *
   * @returns {JobService} The singleton instance of JobService.
   */
  public static getInstance(): JobService {
    if (!JobService.instance) {
      JobService.instance = new JobService();
    }
    return JobService.instance;
  }

  /**
   * Fetches job vacancies with optional pagination and filters.
   *
   * @param {Object} options - The options for fetching job vacancies.
   * @param {number} [options.page=1] - The current page of results.
   * @param {number} [options.pageSize=15] - The number of results per page.
   * @param {{ title?: string; company?: string }} [options.filters={}] - Filters to apply to the job search.
   * @returns {Promise<JobsResponse>} A promise resolving to the fetched job vacancies data.
   * @throws {Error} If the API request fails.
   */
  public async fetchJobVacancies({
    page = 1,
    pageSize = 15,
    filters = {},
  }: {
    page?: number;
    pageSize?: number;
    filters?: { title?: string; company?: string };
  } = {}): Promise<JobsResponse> {
    try {
      const jwtService = JwtService.getInstance();
      const token = await jwtService.getToken();
      const queryParams = this.buildQueryParams(page, pageSize, filters);

      const response = await axiosInstance.get(`/v1/vacancies?${queryParams}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      return { jobs: response.data };
    } catch (error) {
      console.error("Error fetching job vacancies:", error);
      throw new Error("Could not fetch job vacancies.");
    }
  }

  /**
   * Builds query parameters for the API request.
   *
   * @param {number} page - The current page of results.
   * @param {number} pageSize - The number of results per page.
   * @param {{ title?: string; company?: string }} filters - Filters to apply to the job search.
   * @returns {string} A URL-encoded string of query parameters.
   * @private
   */
  private buildQueryParams(
    page: number,
    pageSize: number,
    filters: { title?: string; company?: string },
  ): string {
    const queryParams = new URLSearchParams({
      page: page.toString(),
      page_size: pageSize.toString(),
      ...filters,
    });
    return queryParams.toString();
  }
}

export default JobService;
