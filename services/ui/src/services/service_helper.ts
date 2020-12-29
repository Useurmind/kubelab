export function fetchFromService(method: string, url: string): Promise<Response> {
    return fetch(url, {
        method,
        mode: 'cors',
        cache: 'no-cache',
        // headers: {
        //     "Access-Control-Allow-Origin": "*"
        // }
    } as RequestInit)
}