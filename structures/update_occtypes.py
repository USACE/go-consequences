# import numpy as np
import pandas as pd
# import geopandas as gpd
import os
import json
from pprint import pprint


def prettify_occtypes():
    with open("occtypes.json", "r") as f:
        occtypes = json.load(f)
    
    with open("occtypes_new.json", "w") as out:
        json.dump(occtypes, out, indent=4)

def build_null_df(df):
    output = {
        "damagefunctions": {
            "depth": {
                "source": "Assumed values",
                "damagedriver": "depth",
                "damagefunction": {
                    "xvalues":[0,0],
                    "ydistributions": [
                        {"type": "TriangularDistribution","parameters":{"min":0,"mostlikely":0,"max":0}},
                        {"type": "TriangularDistribution","parameters":{"min":0,"mostlikely":0,"max":0}}
                    ]
                }
            }
        }
    }

    return(output)

def build_damage_function(df):
    output = {
        "damagefunctions": {
            "depth": {
                "source": "Assumed values",
                "damagedriver": "depth",
                "damagefunction": {
                    "xvalues":[],
                    "ydistributions": []
                }
            }
        }
    }
    
    for i, row in df.iterrows():
        tmin = row['t_rebuild_min']
        tmost = row['t_rebuild_mostlikely']
        tmax = row['t_rebuild_max']

        ydist = {
            "type": "TriangularDistribution",
            "parameters":{
                "min": tmin,
                "mostlikely": tmost,
                "max": tmax
            }
        }

        output['damagefunctions']['depth']['damagefunction']['xvalues'].append(row['dmg_pct'])
        output['damagefunctions']['depth']['damagefunction']['ydistributions'].append(ydist)

    return(output)

def read_occtypes():
    with open("occtypes.json", "r") as f:
        occtypes = json.load(f)
    
    for ot in occtypes['occupancytypes'].keys():
        print(ot)
        # print(occtypes['occupancytypes'][ot]['componentdamagefunctions']['structure']['damagefunctions'].keys())
        # print("----")
    
def print_dfs():
    dfs = pd.read_parquet("rowan_2024a_dmg_fns.parquet")
    dfs['co2_cost_pct_sd'] = (dfs['co2_cost_pct_mean'] - dfs['co2_cost_pct_low']) / 1.96
    dfs = dfs[['occtype', 'flood_depth', 'co2_cost_pct_mean', 'co2_cost_pct_sd']]
    dfs['flood_depth'] = dfs['flood_depth'].round(1)
    dfs = dfs[(dfs['flood_depth'] % 1 == 0) & (dfs['flood_depth'] <= 16)]
    print(dfs)


def main():

    res = pd.read_excel("reconstruction_curves.xlsx", sheet_name = "RES")
    com = pd.read_excel("reconstruction_curves.xlsx", sheet_name = "COM")
    ind = pd.read_excel("reconstruction_curves.xlsx", sheet_name = "IND")
    pub = pd.read_excel("reconstruction_curves.xlsx", sheet_name = "PUB")

    with open("occtypes.json", "r") as f:
        occtypes = json.load(f)

    df_res = build_damage_function(res)
    df_com = build_damage_function(com)
    df_ind = build_damage_function(ind)
    df_pub = build_damage_function(pub)
    dfnull = build_null_df(res)

    
    occtypes_out = {"occupancytypes":{}}
    for key, o in occtypes['occupancytypes'].items():
        occtypes_out["occupancytypes"][key] = o
        if(o['name'][0:3] == "RES"):
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['reconstruction'] = df_res
        elif(o['name'][0:3] == "COM"):
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['reconstruction'] = df_com
        elif(o['name'][0:3] == "IND"):
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['reconstruction'] = df_ind
        elif(o['name'][0:3] == "PUB"):
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['reconstruction'] = df_pub
        else:
            occtypes_out['occupancytypes'][o['name']]['componentdamagefunctions']['reconstruction'] = dfnull

    with open("occtypes_reconstruction.json", "w") as out:
        json.dump(occtypes_out, out, indent=4)




if __name__ == "__main__":
    os.chdir(os.path.dirname(os.path.realpath(__file__)))
    main()