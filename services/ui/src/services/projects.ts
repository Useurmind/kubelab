import { IUIConfig } from '../models/config';
import { IGroup } from '../models/project';
import { fetchFromService } from './service_helper';

function fetchFromProjects(config: IUIConfig, method: string, path: string): Promise<Response> {
    return fetchFromService(method, config.projectsBaseUrl + path)
}

export function listGroups(config: IUIConfig): Promise<IGroup[]> {
    return fetchFromProjects(config, "GET", "/groups").then(r => r.json())
        .then(j => j as IGroup[])
}
