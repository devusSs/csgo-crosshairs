export interface userSuccessResponse {
    code: number;
    data: {
        message:string;
        role:string,
        created_at:string,
        e_mail:string,
    };
}

export interface crosshairSuccessResponse{
    code: number;
    data: {
        crosshairs: Crosshair[]
    }
}

export interface User {
    e_mail: string;
    password: string;
}

export interface UserResponse {
  created_at: string,
  e_mail: string,
  role: string
  profile_picture_link: string
}

export interface Email {
    e_mail: string;
}

export interface errorResponse {
    code: number;
    error: {
        error_code: string;
        error_message: string;
    };
}

export interface Crosshair {
    added: string;
    code: string; 
    note: string;
}