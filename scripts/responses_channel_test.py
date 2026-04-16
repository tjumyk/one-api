import argparse
import csv
import json
import sys
import threading
import time
from typing import Any

from openai import OpenAI


def get_output_text(resp: Any) -> str:
    text = getattr(resp, "output_text", None)
    if text:
        return text
    try:
        chunks = []
        for item in getattr(resp, "output", []) or []:
            for content in getattr(item, "content", []) or []:
                if getattr(content, "type", "") in ("output_text", "text"):
                    value = getattr(content, "text", None)
                    if value:
                        chunks.append(value)
        return "\n".join(chunks)
    except Exception:
        return ""


def get_response_json(resp: Any) -> str:
    try:
        if hasattr(resp, "model_dump"):
            payload = resp.model_dump(mode="json")
        elif hasattr(resp, "to_dict"):
            payload = resp.to_dict()
        else:
            payload = json.loads(resp.model_dump_json())
        return json.dumps(payload, ensure_ascii=False, separators=(",", ":"))
    except Exception:
        return ""


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--base-url", required=True)
    parser.add_argument("--api-key", required=True)
    parser.add_argument("--channel-id", required=True)
    parser.add_argument("--model", required=True)
    parser.add_argument("--claimed-model", required=True)
    parser.add_argument("--max-output-tokens", type=int, default=128)
    parser.add_argument("--csv", required=False)
    args = parser.parse_args()

    api_key = f"{args.api_key}-{args.channel_id}"
    client = OpenAI(base_url=args.base_url.rstrip("/") + "/v1", api_key=api_key)

    start = time.time()
    row = {
        "channel_id": args.channel_id,
        "requested_model": args.model,
        "claimed_model": args.claimed_model,
        "response_model": "",
        "max_output_tokens": args.max_output_tokens,
        "success": False,
        "match": False,
        "latency_ms": 0,
        "error": "",
        "response_json": "",
        "response_message": "",
        "response_message_preview": "",
    }

    stop_event = threading.Event()

    def keepalive() -> None:
        tick = 0
        while not stop_event.wait(5):
            tick += 1
            print(
                json.dumps(
                    {
                        "type": "keepalive",
                        "channel_id": args.channel_id,
                        "model": args.model,
                        "claimed_model": args.claimed_model,
                        "elapsed_s": tick * 5,
                    },
                    ensure_ascii=True,
                ),
                flush=True,
            )

    keeper = threading.Thread(target=keepalive, daemon=True)
    keeper.start()

    try:
        resp = client.responses.create(
            model=args.model,
            input="<What model are you? what is your knowledge cutoff?>",
            max_output_tokens=args.max_output_tokens,
        )
    finally:
        stop_event.set()
        keeper.join(timeout=1)

    try:
        latency_ms = int((time.time() - start) * 1000)
        latency_ms = int((time.time() - start) * 1000)
        response_model = getattr(resp, "model", "") or ""
        response_json = get_response_json(resp)
        response_message = get_output_text(resp)
        response_message_preview = response_message[:200]
        row.update({
            "response_model": response_model,
            "success": True,
            "match": response_model == args.claimed_model,
            "latency_ms": latency_ms,
            "response_json": response_json,
            "response_message": response_message,
            "response_message_preview": response_message_preview,
        })
    except Exception as e:
        row.update({
            "success": False,
            "match": False,
            "latency_ms": int((time.time() - start) * 1000),
            "error": str(e),
        })

    if args.csv:
        with open(args.csv, "a", newline="", encoding="utf-8") as f:
            writer = csv.DictWriter(
                f,
                fieldnames=[
                    "channel_id",
                    "requested_model",
                    "claimed_model",
                    "response_model",
                    "max_output_tokens",
                    "success",
                    "match",
                    "latency_ms",
                    "error",
                    "response_json",
                    "response_message",
                    "response_message_preview",
                ],
            )
            if f.tell() == 0:
                writer.writeheader()
            writer.writerow({k: row.get(k, "") for k in [
                "channel_id",
                "requested_model",
                "claimed_model",
                "response_model",
                "max_output_tokens",
                "success",
                "match",
                "latency_ms",
                "error",
                "response_json",
                "response_message",
                "response_message_preview",
            ]})

    print(json.dumps(row, ensure_ascii=True))
    return 0 if row["success"] and row["match"] else 1


if __name__ == "__main__":
    raise SystemExit(main())
