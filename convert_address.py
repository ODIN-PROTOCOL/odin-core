import sys
import bech32

def convert_chain_address(prefix, address_string):
    _, data = bech32.bech32_decode(address_string)
    return bech32.bech32_encode(prefix, data)

if __name__ == '__main__':
    if len(sys.argv) != 3:
        print("Usage: python convert_address.py <prefix> <address>")
        sys.exit(1)
    
    prefix = sys.argv[1]
    address = sys.argv[2]
    
    print(convert_chain_address(prefix, address))