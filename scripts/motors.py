from argparse import ArgumentParser
from flask import Flask

def init_eh():
    try:
        return __import__("explorerhat")
    except ImportError:
        print("Explorer Hat module is not present")
        return None


app = Flask(__name__)
eh = init_eh()


def forward():
    eh.motor.one.backwards(100)
    eh.motor.two.forwards(100)


def backward():
    eh.motor.one.forwards(100)
    eh.motor.two.backwards(100)


def left():
    eh.motor.two.stop()
    eh.motor.one.backwards(100)


def right():
    eh.motor.one.stop()
    eh.motor.two.forwards(100)


def stop():
    eh.motor.one.stop()
    eh.motor.two.stop()


def clockwise():
    eh.motor.one.forwards(100)
    eh.motor.two.forwards(100)


def anti_clockwise():
    eh.motor.one.backwards(100)
    eh.motor.two.backwards(100)


stateToFunction = {
    "forward": forward,
    "back": backward,
    "left": left,
    "right": right,
    "stop": stop,
    "clockwise": clockwise,
    "anti-clockwise": anti_clockwise
}


def update_state(state):
    try:
        f = stateToFunction[state]
        if eh is None:
            print("State #" + f.__name__ + "# was called")
        else:
            f()
    except KeyError:
        print("No state with name " + state)


@app.route("/<state>")
def update_robot(state=None):
    update_state(state)
    return "ok"


if __name__ == "__main__":
    parser = ArgumentParser(description='Controller for Explorer Hat Motors')
    parser.add_argument('-p', '--port', help='Port on which the http server will run')

    args = parser.parse_args()
    app.run(host='0.0.0.0', port=80)
