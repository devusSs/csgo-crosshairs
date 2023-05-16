import BigNumber from 'bignumber.js';

// Intentionally no 0 and 1 number in DICTIONARY
const DICTIONARY = 'ABCDEFGHJKLMNOPQRSTUVWXYZabcdefhijkmnopqrstuvwxyz23456789';
const DICTIONARY_LENGTH = DICTIONARY.length;
const SHARECODE_PATTERN = /CSGO(-?[\w]{5}){5}$/;

interface CrosshairData {
  cl_crosshairgap: number;
  cl_crosshair_outlinethickness: number;
  cl_crosshaircolor_r: number;
  cl_crosshaircolor_g: number;
  cl_crosshaircolor_b: number;
  cl_crosshairalpha: number;
  cl_crosshair_dynamic_splitdist: number;
  cl_fixedcrosshairgap: number;
  cl_crosshaircolor: number;
  cl_crosshair_drawoutline: number;
  cl_crosshair_dynamic_splitalpha_innermod: number;
  cl_crosshair_dynamic_splitalpha_outermod: number;
  cl_crosshair_dynamic_maxdist_splitratio: number;
  cl_crosshairthickness: number;
  cl_crosshairstyle: number;
  cl_crosshairdot: number;
  cl_crosshairgap_useweaponvalue: number;
  cl_crosshairusealpha: number;
  cl_crosshair_t: number;
  cl_crosshairsize: number;
}

function bigNumberToByteArray(big: BigNumber): number[] {
  const str : string = big.toString(16).padStart(36, '0');
  const bytes : number[] = [];
  for (let i = 0; i < str.length; i += 2) {
    const end = i+2;
    bytes.push(parseInt(str.slice(i, end), 16));
  }

  return bytes;
}


function parseBytes(bytes: number[]): CrosshairData {
    return {
      cl_crosshairgap: Int8Array.of(bytes[2])[0] / 10.0,
  
      cl_crosshair_outlinethickness: (bytes[3] & 7) / 2.0,
  
      cl_crosshaircolor_r: bytes[4],
      cl_crosshaircolor_g: bytes[5],
      cl_crosshaircolor_b: bytes[6],
      cl_crosshairalpha: bytes[7],
      cl_crosshair_dynamic_splitdist: bytes[8],
  
      cl_fixedcrosshairgap: Int8Array.of(bytes[9])[0] / 10.0,
  
      cl_crosshaircolor: bytes[10] & 7,
      cl_crosshair_drawoutline: bytes[10] & 8 ? 1 : 0,
      cl_crosshair_dynamic_splitalpha_innermod: ((bytes[10] & 0xF0) >> 4) / 10.0,
  
      cl_crosshair_dynamic_splitalpha_outermod: (bytes[11] & 0xF) / 10.0,
      cl_crosshair_dynamic_maxdist_splitratio: ((bytes[11] & 0xF0) >> 4) / 10.0,
  
      cl_crosshairthickness: (bytes[12] & 0x3F) / 10.0,
  
      cl_crosshairstyle: (bytes[13] & 0xE) >> 1,
      cl_crosshairdot: bytes[13] & 0x10 ? 1 : 0,
      cl_crosshairgap_useweaponvalue: bytes[13] & 0x20 ? 1 : 0,
      cl_crosshairusealpha: bytes[13] & 0x40 ? 1 : 0,
      cl_crosshair_t: bytes[13] & 0x80 ? 1 : 0,
  
      cl_crosshairsize: (((bytes[15] & 0x1f) << 8) + bytes[14]) / 10.0
    };
  }
  
  export const decode = shareCode => {
    if (!shareCode.match(SHARECODE_PATTERN)) {
      throw new Error('Invalid share code');
    }
  
    shareCode = shareCode.replace(/CSGO|-/g, '');
    const chars : any = Array.from(shareCode).reverse();
    let big = new BigNumber(0);
  
    for (let i = 0; i < chars.length; i++) {
      big = big.multipliedBy(DICTIONARY_LENGTH).plus(DICTIONARY.indexOf(chars[i]));
    }
    
    return parseBytes(bigNumberToByteArray(big));
  }
  