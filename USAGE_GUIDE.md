# Intent Recognition Usage Guide

## Overview
This system uses natural language processing to understand your intent and extract relevant information. For the best results, we recommend using **quotes** for titles and names.

## Recommended Usage

### ✅ **Best Practice: Use Quotes**

#### Calendar Events
```bash
# Good - Clear and unambiguous
create calendar event "team meeting" tomorrow at 2pm
schedule appointment "doctor visit" next week
add meeting "project review" to calendar for Friday
```

#### Tasks
```bash
# Good - Clear and unambiguous
create new task "buy groceries" for tomorrow
add todo item "prepare presentation" by Friday
create reminder task "call dentist" next week
```

#### Contacts
```bash
# Good - Clear and unambiguous
create a new contact named "Alice Smith" with email alice@example.com
add contact "John Doe" with phone 555-123-4567
save contact "Bob Wilson" with email bob@example.com
```

### ⚠️ **Fallback: Without Quotes**
The system will try its best to understand your intent even without quotes, but results may be less accurate:

```bash
# Works but may be less accurate
create calendar event for team meeting tomorrow at 2pm
create new task called buy groceries
add contact John Smith with email john@example.com
```

## What Gets Extracted

### Events
- **Title**: The name of the event
- **Date**: When the event occurs (tomorrow, next week, etc.)
- **Time**: What time the event is scheduled
- **Location**: Where the event takes place (if specified)

### Tasks
- **Title**: The name of the task
- **Due Date**: When the task is due (if specified)
- **Priority**: Task priority (if specified)

### Contacts
- **Name**: The person's name
- **Email**: Email address (if provided)
- **Phone**: Phone number (if provided)

## Examples

### Complete Examples
```bash
# Event with all details
create calendar event "team standup" tomorrow at 9am in conference room

# Task with due date
create new task "prepare quarterly report" due next Friday

# Contact with full details
create contact "Sarah Johnson" with email sarah@company.com and phone 555-987-6543
```

### Partial Examples (System will ask for missing info)
```bash
# Missing title
create calendar event tomorrow at 2pm
# System asks: "What should I call this event? (use quotes for clarity)"

# Missing time
create calendar event "team meeting" tomorrow
# System asks: "What time should this event be?"

# Missing name
create a new contact with email alice@example.com
# System asks: "What's the name of the contact?"
```

## Tips for Best Results

1. **Use quotes for titles and names** - This makes the system's job much easier
2. **Be specific about dates and times** - Use clear terms like "tomorrow", "next week", "2pm"
3. **Include all relevant details** - The more information you provide, the better
4. **Use natural language** - The system understands conversational language

## System Capabilities

The system can recognize these intents:
- **CreateEvent**: Calendar events and meetings
- **CreateTask**: Tasks and todo items
- **CreateContact**: Adding new contacts
- **Weather**: Weather information requests
- **Time**: Current time requests
- **Calculator**: Mathematical calculations
- And more...

Each intent has specific required fields and will prompt you for missing information. 