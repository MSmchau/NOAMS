"""
Netmiko SSH Automation Worker

A lightweight microservice that executes SSH commands on network devices
using Netmiko. Supports H3C Comware, Huawei VRP, Cisco IOS, Ruijie OS, etc.
"""

from flask import Flask, request, jsonify
from netmiko import ConnectHandler, NetmikoTimeoutException, NetmikoAuthenticationException
import logging
import time

app = Flask(__name__)
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("netmiko-worker")

DEVICE_TYPE_MAP = {
    "hp_comware": "hp_comware",
    "huawei": "huawei",
    "cisco_ios": "cisco_ios",
    "ruijie_os": "ruijie_os",
}


def build_device_params(data):
    device_type = DEVICE_TYPE_MAP.get(data.get("device_type"))
    if not device_type:
        raise ValueError(f"Unsupported device type: {data.get('device_type')}")

    params = {
        "device_type": device_type,
        "host": data["ip"],
        "port": data.get("port", 22),
        "username": data["username"],
        "password": data["password"],
        "timeout": data.get("timeout", 30),
        "global_delay_factor": data.get("delay_factor", 1),
    }

    if data.get("secret"):
        params["secret"] = data["secret"]

    return params


@app.route("/api/v1/execute", methods=["POST"])
def execute_commands():
    """Execute one or more commands on a network device."""
    data = request.get_json()
    if not data:
        return jsonify({"success": False, "error": "No JSON data provided"}), 400

    required = ["ip", "username", "password", "device_type", "commands"]
    for field in required:
        if field not in data:
            return jsonify({"success": False, "error": f"Missing field: {field}"}), 400

    try:
        device_params = build_device_params(data)
    except ValueError as e:
        return jsonify({"success": False, "error": str(e)}), 400

    results = []
    start_time = time.time()

    try:
        logger.info(f"Connecting to {data['ip']} ({data['device_type']})")
        connection = ConnectHandler(**device_params)

        if data.get("enable_mode", False):
            connection.enable()

        for cmd in data["commands"]:
            try:
                logger.info(f"  Executing: {cmd}")
                output = connection.send_command(
                    cmd,
                    delay_factor=data.get("cmd_delay_factor", 1),
                    expect_string=data.get("expect_string"),
                )
                results.append({
                    "command": cmd,
                    "output": output,
                    "success": True,
                })
            except Exception as e:
                logger.error(f"  Command failed: {cmd} - {str(e)}")
                results.append({
                    "command": cmd,
                    "output": "",
                    "success": False,
                    "error": str(e),
                })

        connection.disconnect()

        elapsed = time.time() - start_time
        logger.info(f"Completed {data['ip']} in {elapsed:.1f}s")

        return jsonify({
            "success": True,
            "results": results,
            "elapsed": round(elapsed, 2),
        })

    except NetmikoTimeoutException as e:
        logger.error(f"Timeout connecting to {data['ip']}: {str(e)}")
        return jsonify({"success": False, "error": f"Connection timeout: {str(e)}"}), 504

    except NetmikoAuthenticationException as e:
        logger.error(f"Auth failed for {data['ip']}: {str(e)}")
        return jsonify({"success": False, "error": "Authentication failed"}), 401

    except Exception as e:
        logger.error(f"Error connecting to {data['ip']}: {str(e)}")
        return jsonify({"success": False, "error": str(e)}), 500


@app.route("/api/v1/health", methods=["GET"])
def health():
    return jsonify({"status": "ok"})


@app.route("/api/v1/test", methods=["POST"])
def test_connection():
    """Test SSH connectivity to a device without executing commands."""
    data = request.get_json()
    if not data:
        return jsonify({"success": False, "error": "No JSON data provided"}), 400

    try:
        test_params = build_device_params(data)
        test_params["commands"] = ["display clock", "display version | include Version"]
    except ValueError as e:
        return jsonify({"success": False, "error": str(e)}), 400

    try:
        logger.info(f"Testing connection to {data['ip']}")
        connection = ConnectHandler(**test_params)
        output = connection.send_command("display clock", delay_factor=1)
        connection.disconnect()
        return jsonify({"success": True, "output": output})
    except Exception as e:
        return jsonify({"success": False, "error": str(e)}), 500


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000, debug=False)
