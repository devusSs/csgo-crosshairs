const DICTIONARY = "ABCDEFGHJKLMNOPQRSTUVWXYZabcdefhijkmnopqrstuvwxyz23456789";
const DICTIONARY_LENGTH = BigInt(DICTIONARY.length);

type Info = {
  cl_crosshairgap: number;
  cl_crosshair_outlinethickness: number;
  cl_crosshaircolor_r: number;
  cl_crosshaircolor_g: number;
  cl_crosshaircolor_b: number;
  cl_crosshairalpha: number;
  cl_crosshair_dynamic_splitdist: number;
  cl_fixedcrosshairgap: number;
  cl_crosshaircolor: number;
  cl_crosshair_drawoutline: boolean;
  cl_crosshair_dynamic_splitalpha_innermod: number;
  cl_crosshair_dynamic_splitalpha_outermod: number;
  cl_crosshair_dynamic_maxdist_splitratio: number;
  cl_crosshairthickness: number;
  cl_crosshairstyle: number;
  cl_crosshairdot: boolean;
  cl_crosshairgap_useweaponvalue: boolean;
  cl_crosshairusealpha: boolean;
  cl_crosshair_t: boolean;
  cl_crosshairsize: number;
};



const serializeToBytes = (info : Info): number[] => {
    const bytes = [
        0,
        1,
        (info.cl_crosshairgap * 10) & 0xff,
        (info.cl_crosshair_outlinethickness * 2) & 7,
        info.cl_crosshaircolor_r,
        info.cl_crosshaircolor_g,
        info.cl_crosshaircolor_b,
        info.cl_crosshairalpha,
        info.cl_crosshair_dynamic_splitdist,
        (info.cl_fixedcrosshairgap * 10) & 0xff,
        (info.cl_crosshaircolor & 7) |
            (info.cl_crosshair_drawoutline ? 8 : 0) |
            (info.cl_crosshair_dynamic_splitalpha_innermod * 10) << 4,
        ((info.cl_crosshair_dynamic_splitalpha_outermod * 10) & 0xf) |
            ((info.cl_crosshair_dynamic_maxdist_splitratio * 10) << 4),
        (info.cl_crosshairthickness * 10) & 0x3f,
        ((info.cl_crosshairstyle << 1) & 0xe) |
            (info.cl_crosshairdot ? 0x10 : 0) |
            (info.cl_crosshairgap_useweaponvalue ? 0x20 : 0) |
            (info.cl_crosshairusealpha ? 0x40 : 0) |
            (info.cl_crosshair_t ? 0x80 : 0),
        (info.cl_crosshairsize * 10) & 0xff,
        ((info.cl_crosshairsize * 10) >> 8) & 0x1f,
        0,
        0
    ];

    const sum = bytes.slice(1).reduce((acc, byte) => acc + byte, 0);
    bytes[0] = sum & 0xff;

    return bytes;
};

export const encode = (info: Info): string => {
    const bytes = serializeToBytes(info);
    
    let acc = BigInt(0);
    let pos = BigInt(1);
    for (let i = bytes.length - 1; i >= 0; i--) {
      acc += BigInt(bytes[i]) * pos;
      pos *= 256n;
    }
    
    let result = '';
    for (let i = 0; i < 25; ++i) {
        const digit = acc % DICTIONARY_LENGTH;
        acc = acc / DICTIONARY_LENGTH;
        result += DICTIONARY.charAt(Number(digit));
    }
    
    return `CSGO-${result.slice(0, 5)}-${result.slice(5, 10)}-${result.slice(10, 15)}-${result.slice(15, 20)}-${result.slice(20, 25)}`;
};
