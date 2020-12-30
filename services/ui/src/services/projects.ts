import { IUIConfig } from '../models/config';
import { IGroup } from '../models/project';
import { fetchFromService } from './service_helper';

function fetchFromProjects(config: IUIConfig, method: string, path: string, body: any): Promise<Response> {
    return fetchFromService(method, config.projectsBaseUrl + path, body)
}

export function createGroup(config: IUIConfig, group: IGroup): Promise<IGroup> {
    return fetchFromProjects(config, "POST", "/groups", group).then(r => r.json())
        .then(j => j as IGroup)
}

export function deleteGroup(config: IUIConfig, groupId: number): Promise<boolean> {
    return fetchFromProjects(config, "DELETE", `/groups/${groupId}`, null).then(r => r.status === 204)
}

export function listGroups(config: IUIConfig): Promise<IGroup[]> {
    return fetchFromProjects(config, "GET", "/groups", null).then(r => r.json())
        .then(j => j as IGroup[])
}
