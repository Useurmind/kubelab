export function fetchFromService(method: string, url: string, body: any): Promise<Response> {
    return fetch(url, {
        method,
        mode: 'cors',
        cache: 'no-cache',
        body: body ? JSON.stringify(body) : null
        // headers: {
        //     "Access-Control-Allow-Origin": "*"
        // }
    } as RequestInit)
}