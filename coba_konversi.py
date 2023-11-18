import math
import numpy as np

def calcx(az, el, h):
    B_rad = np.radians(90 - az)
    C_rad = np.radians(el)

    # Calculate the result
    x = np.cos(B_rad) * (h / np.tan(C_rad))
#    x1 = np.cos(np.radians((90-az) * (np.pi * np.radians(180)) * (h/(np.tan(np.radians((np.pi / np.radians(180)) / el))))  ))
    return x

def calcy(az, el, h):
    B_rad = np.radians(90 - az)
    C_rad = np.radians(el)

    # Calculate the result
    y = np.sin(B_rad) * (h / np.tan(C_rad))
    return y

def calc_windspeed(x1, x2, y1, y2):
    result = np.sqrt((x2 - x1)**2 + (y2 - y1)**2) / 120
    return result

def calc_winddir(x1, x2, y1, y2):
    condition = (x2 - x1) > 0
    dir = np.where(condition, (4.712 - np.arctan((y2 - y1) / (x2 - x1))) / (np.pi / 180),
                      (4.712 - np.arctan((y2 - y1) / (x2 - x1)) - 3.1416) / (np.pi / 180))

    return dir

def calc_speed_dir(elevation1, azimuth1, elevation2, azimuth2, h1, h2):
#    #find x and y pos
#    #COS((90-B11)*(PI()/180))*(K11/TAN ((PI()/180)*C11))
#    #SIN((PI()/180) (90-B11))*(K11/TAN((PI()/180)*C11))
    x1 = calcx(azimuth1, elevation1, h1)
    x2 = calcx(azimuth2, elevation2, h2)
    y1 = calcy(azimuth1, elevation1, h1)
    y2 = calcy(azimuth2, elevation2, h2)

    #find windspeed
    speed = calc_windspeed(x1, x2, y1, y2)
    dir = calc_winddir(x1, x2, y1, y2)
    return speed, dir

#    
#    return x1
# Example usage
elevation1 = 20  # replace with your actual elevation angle for point 1
azimuth1 = 158  # replace with your actual azimuth angle for point 1
elevation2 = 23.2  # replace with your actual elevation angle for point 2
azimuth2 = 160.8  # replace with your actual azimuth angle for point 2

az1 = 230.3
az2 = 223
el1 = 28.2
el2 = 27.9
h1 = 153
h2 = 433
#x1 = calcx(az, el, h)
#y1 = calcy(az, el, h)
#print('cek rumus', x1, y1)
speed, dir = calc_speed_dir(el1, az1, el2, az2, h1, h2)
print('cek hasil', speed, dir)