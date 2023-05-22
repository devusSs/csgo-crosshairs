export interface successResponse {
    code: number;
    data: {};
}

export interface User {
    e_mail: string;
    password: string;
}

export interface errorResponse {
    code: number;
    error: {
        error_code: string;
        error_message: string;
    };
}

export interface crosshair {
    added: string;
    code: string; 
    note: string;
}