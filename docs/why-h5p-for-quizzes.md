# WHy H5P for Quizzes?

H5P has become a popular way to represent quiz and assessment information, and here’s why it’s a good option for your use case:

---

### **1. Structured but Lightweight**

* H5P packages are just **zip files** containing JSON (`content.json`, `semantics.json`) and assets.
* Much simpler to handle than XML-based standards like QTI.
* Easy to parse in any modern language (Go, Python, JS, etc.) and store in databases.

---

### **2. Flexible Question Types**

* Out of the box, H5P supports many **quiz and interaction types**: multiple choice, fill-in-the-blank, drag-and-drop, true/false, summary, interactive video quizzes, etc.
* Each type is defined by a `semantics.json` schema, so the data format is consistent and validated.

---

### **3. Open and Extensible**

* Fully **open-source** under MIT license.
* Developers can create new content types by defining their own `semantics.json`.
* Easy to extend for custom question types or metadata you may need in a question bank.

---

### **4. Widely Supported**

* Works in popular LMS/CMS like Moodle, WordPress, Drupal.
* Many e-learning platforms already support H5P import/export.
* Good community and ecosystem (compared to rolling your own schema).

---

### **5. Human-Friendly JSON**

* Example from H5P Multiple Choice (`content.json`):

```json
{
  "question": "What is 2 + 2?",
  "answers": [
    { "text": "3", "correct": false },
    { "text": "4", "correct": true },
    { "text": "5", "correct": false }
  ],
  "behaviour": {
    "singleAnswer": true,
    "randomAnswers": true
  }
}
```

This is straightforward to read, write, and validate — much cleaner than QTI XML.

---

### **6. Future-Proof with Interoperability**

* You can start with H5P JSON for **ease of use**.
* Later, you can **map to QTI** if you need LMS-grade interoperability (since question/answer/feedback concepts map pretty well).
* You can also integrate with **xAPI** for tracking learner interactions.

---

✅ **In short:**
H5P is a great option because it balances **developer-friendliness (JSON, open source)** with **educational robustness (rich quiz types, LMS support)**. It’s an approachable starting point for a question bank schema, while still leaving the door open to more formal standards like QTI if needed.
