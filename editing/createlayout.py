from flask import Flask, render_template, Markup, request, redirect, url_for
import firebase_admin
from firebase_admin import credentials
from firebase_admin import firestore
import random

# Use a service account
cred = credentials.Certificate('/path-to-service-account.json') #PROD
# cred = credentials.Certificate('dev.json') #DEV

firebase_admin.initialize_app(cred)

db = firestore.client()

# Update to initialize db for new species
species_letter = "e"
story_collection = "story"
pattern = [2, 3, 2, 2, 3]

zero_data = {"t": "Main text", "a": {}}
for z in range(pattern[0]):
    zero_data["a"][f"a{z+1}"] = {"t": "Answer text", "x": f"{species_letter}{z+1}"}
db.collection(story_collection).document(f"{species_letter}").set(zero_data)

for a in range(pattern[0]):
    data_a = {"t": "Main text", "a": {}}
    for v in range(pattern[1]):
        data_a["a"][f"a{v+1}"] = {"t": "Answer text", "x": f"{species_letter}{a+1}{v+1}"}
    db.collection(story_collection).document(f"{species_letter}{a+1}").set(data_a)

    for b in range(pattern[1]):
        data_b = {"t": "Main text","a": {}}
        for w in range(pattern[2]):
            data_b["a"][f"a{w+1}"] = {"t": "Answer text", "x": f"{species_letter}{a+1}{b+1}{w+1}"}
        db.collection(story_collection).document(f"{species_letter}{a+1}{b+1}").set(data_b)

        for c in range(pattern[2]):
            data_c = {"t": "Main text","a": {}}
            for x in range(pattern[3]):
                data_c["a"][f"a{x+1}"] = {"t": "Answer text", "x": f"{species_letter}{a+1}{b+1}{c+1}{x+1}"}
            db.collection(story_collection).document(f"{species_letter}{a+1}{b+1}{c+1}").set(data_c)

            for d in range(pattern[3]):
                data_d = {"t": "Main text","a": {}}
                for y in range(pattern[4]):
                    data_d["a"][f"a{y+1}"] = {"t": "Answer text", "x": f"{species_letter}{a+1}{b+1}{c+1}{d+1}{y+1}"}
                db.collection(story_collection).document(f"{species_letter}{a+1}{b+1}{c+1}{d+1}").set(data_d)

                for e in range(pattern[4]):
                    db.collection(story_collection).document(f"{species_letter}{a+1}{b+1}{c+1}{d+1}{e+1}").set({
                        "o": "Main text",
                        "n": "NAMEEARNED",
                        "v": random.randint(2, 10)
                    })

print("Done")

