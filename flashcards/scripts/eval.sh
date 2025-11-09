#!/bin/bash

# Output file
OUTPUT_FILE="scripts/eval_output.txt"

# Clear the output file
> "$OUTPUT_FILE"

# Function to call the quiz API
call_quiz_api() {
  local payload="$1"
  local scenario_name="$2"

  echo "==================================================" >> "$OUTPUT_FILE"
  echo "Scenario: $scenario_name" >> "$OUTPUT_FILE"
  echo "==================================================" >> "$OUTPUT_FILE"
  echo "Input Payload:" >> "$OUTPUT_FILE"
  echo "$payload" >> "$OUTPUT_FILE"
  echo "" >> "$OUTPUT_FILE"
  echo "API Output:" >> "$OUTPUT_FILE"

  curl -s -X POST http://localhost:8080/quiz \
  -H "Content-Type: application/json" \
  -d "$payload" | python -m json.tool >> "$OUTPUT_FILE"

  echo "" >> "$OUTPUT_FILE"
  echo "--------------------------------------------------" >> "$OUTPUT_FILE"
  echo "" >> "$OUTPUT_FILE"
}

# --- Scenarios ---

# Scenario 1: Initial Turn (no user input)
SCENARIO_1_PAYLOAD='{
  "messages": []
}'
call_quiz_api "$SCENARIO_1_PAYLOAD" "Initial Turn"

# Scenario 2: Correct Answer
SCENARIO_2_PAYLOAD='{
  "messages": [
    {
      "role": "assistant",
      "content": "The Kingdom of England is traditionally dated from the rule of Æthelstan from what year?"
    },
    {
      "role": "user",
      "content": "927"
    }
  ]
}'
call_quiz_api "$SCENARIO_2_PAYLOAD" "Correct Answer"

# Scenario 3: Incorrect Answer
SCENARIO_3_PAYLOAD='{
  "messages": [
    {
      "role": "assistant",
      "content": "In 1066, the Norman Conquest was led by whom?"
    },
    {
      "role": "user",
      "content": "King Arthur"
    }
  ]
}'
call_quiz_api "$SCENARIO_3_PAYLOAD" "Incorrect Answer"

# Scenario 4: Off-Topic Question
SCENARIO_4_PAYLOAD='{
  "messages": [
    {
      "role": "assistant",
      "content": "The Magna Carta was signed by which king in 1215?"
    },
    {
      "role": "user",
      "content": "What is the weather like in London?"
    }
  ]
}'
call_quiz_api "$SCENARIO_4_PAYLOAD" "Off-Topic Question"

# Scenario 5: Gibberish Input
SCENARIO_5_PAYLOAD='{
  "messages": [
    {
      "role": "assistant",
      "content": "The Hundred Years War was a series of conflicts between England and which other kingdom?"
    },
    {
      "role": "user",
      "content": "asdfghjkl"
    }
  ]
}'
call_quiz_api "$SCENARIO_5_PAYLOAD" "Gibberish Input"

echo "Evaluation script finished. Output written to $OUTPUT_FILE"
