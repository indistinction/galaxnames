from flask import Flask, render_template, Markup, request, redirect, url_for
import firebase_admin
from firebase_admin import credentials
from firebase_admin import firestore

# Use a service account
cred = credentials.Certificate('/path-to-service-account.json')
firebase_admin.initialize_app(cred)

db = firestore.client()

app = Flask(__name__)

@app.route('/')
@app.route('/play/')
def start():
    html = "<a href='/play/c'>Click to begin the Cylindurus storylines.</a>"
    html_m = Markup(html)
    return render_template('index.html', html=html_m)


@app.route('/play/<level>')
def play(level, locked=False):
    html = ""
    if locked:
        html += "<p>Cannot edit right now, this level is being edited by somebody else.</p>"
    html += f'<a href="/play/{level[:-1]}">&lt; previous step</a><br>'
    q = db.collection('story').document(level).get().to_dict()
    if "o" in q:
        # o for outcome
        text = q['o'].replace('\n', '<br>')
        html += f"<p>{text}</p>"
        html += f"""<p>Whenever a new Galaxiator is 'recruited', it is the privilege of the Hunter to provide them
                    with a gladiatorial name for the Galaxiators arena. From today onwards you will be known as
                    <strong>{q['n']}</strong></p>"""
    else:
        # t for text
        text = q['t'].replace('\n', '<br>')
        html += f"<p>{text}</p>"
        for a in q['a']:
            html += f"<a href='/play/{q['a'][a]['x']}' class='btn'>{q['a'][a]['t']}</a><br>"

    html_m = Markup(html + f"<a href='/edit/{level}' class='edit'>Click to edit this text...</a><br><small><em>Tracking code for Gompy: {level}</em></small>")
    return render_template('index.html', html=html_m)


@app.route('/edit/<level>', methods=["GET"])
def edit(level):
    # Need back button that performs unlock first
    q = db.collection('story').document(level).get().to_dict()
    if "lock" in q and q["lock"]:
        return play(level, True)
    else:
        db.collection('story').document(level).update({"lock": True})

    html = f"<form method='POST' action='/save/{level}'>"

    if "o" in q:
        # o for outcome
        html += "<p>Text<br><textarea rows='15' name='o'>" + q['o'] + "</textarea><br>"
        html += f"Name outcome: <input name='n' value='{q['n']}'> $GLXR value <input name='v' value='{q['v'] if q['v'] else 0}'></p>"
    else:
        # t for text
        html += "<p>Text<br><textarea rows='15' name='t'>" + q['t'] + "</textarea><p>"
        i = 0
        for a in q['a']:
            html += f"<p>Choice {i+1}<br>"
            html += f"<textarea rows='8' name='{a}'>" + q['a'][a]['t'] + "</textarea><br></p>"
            i += 1

    html_m = Markup(html+f"<input type='submit' value='Save' class='edit'></form><br><small><em>Tracking code for Gompy: {level}</em></small>")
    return render_template('index.html', html=html_m)


@app.route('/save/<level>', methods=["POST"])
def save(level):
    request.get_data()
    if "o" in request.form:
        db.collection('story').document(level).update({
            "lock": False,
            "n": request.form["n"],
            "v": int(request.form["v"]),
            "o": request.form["o"]
        })
    else:
        a = {}
        for key in request.form:
            if key.startswith("a"):
                a[key] = {
                    "t": request.form[key],
                    "x": level + key.split("a")[1]
                }
        db.collection('story').document(level).update({
            "a": a,
            "t": request.form["t"],
            "lock": False
        })
    return redirect(url_for('play', level=level))


if __name__ == "__main__":
    app.run()

