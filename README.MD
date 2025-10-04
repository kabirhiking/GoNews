# 📰 GoNews - Full-Text News Search Engine

A high-performance, intelligent search engine built in Go that indexes and searches through thousands of news articles using TF-IDF scoring and Boolean query logic.

## 🎯 Project Overview

**GoNews** is a full-text search engine that demonstrates core information retrieval concepts including:
- Inverted index construction
- TF-IDF (Term Frequency-Inverse Document Frequency) scoring
- Boolean query processing (AND, OR, NOT)
- Phrase search with exact matching
- Text tokenization and optional stemming
- Real-time search with millisecond response times

## 🚀 Features

✅ **Fast Indexing**: Index 6,000+ documents in ~3 seconds  
✅ **Lightning Search**: Execute queries in ~1 millisecond  
✅ **Boolean Logic**: Support for AND, OR, NOT operators  
✅ **Phrase Search**: Exact phrase matching with `"quoted text"`  
✅ **Complex Queries**: Nested queries with parentheses  
✅ **TF-IDF Scoring**: Intelligent relevance ranking  
✅ **Snippet Generation**: Contextual preview of matches  
✅ **Optional Stemming**: Root word matching (running → run)  

## 📁 Project Structure

```
GoNews/
├── main.go          # Entry point - CLI interface, timing, orchestration
├── load.go          # CSV loader - parses news dataset
├── analyze.go       # Text analysis - tokenization, stemming, normalization
├── index.go         # Inverted index - document indexing & search engine
├── query.go         # Query processor - RPN parsing, Boolean logic, snippets
├── go.mod           # Go module dependencies
└── GoNews/
    ├── README.md    # This file
    └── data/
        ├── news.csv         # Main news dataset (6,335 articles)
        ├── Fake.csv         # Fake news articles
        ├── True.csv         # Real news articles
        ├── clean_news.csv   # Processed dataset
        └── new.ipynb        # Python notebook for data download
```

## 🛠️ Installation & Setup

### Prerequisites
- Go 1.16+ installed
- Python 3.x (for data download)
- Git

### 1. Clone or Navigate to Project
```bash
cd C:\Users\pc\Desktop\GoNews
```

### 2. Initialize Go Module (if needed)
```bash
go mod init news-search
go mod tidy
```

### 3. Download Dataset (Optional - already included)
The dataset is already included in `GoNews/data/news.csv`. If you want to download fresh data:

```python
# In new.ipynb or Python script
import kagglehub
path = kagglehub.dataset_download("clmentbisaillon/fake-and-real-news-dataset")
# Copy CSV files to GoNews/data/
```

## 🎮 Usage

### Basic Search
```bash
# Search for articles containing "climate" OR "change"
go run . -p GoNews/data/news.csv -q "climate change" -n 10
```

### Phrase Search (Exact Match)
```bash
# Search for exact phrase "climate change"
go run . -p GoNews/data/news.csv -q "\"climate change\"" -n 5
```

### Boolean Queries

#### AND Operator
```bash
# Both terms must be present
go run . -p GoNews/data/news.csv -q "climate AND change" -n 5
```

#### OR Operator
```bash
# Either term can be present
go run . -p GoNews/data/news.csv -q "climate OR environment" -n 5
```

#### NOT Operator
```bash
# Exclude documents with "hoax"
go run . -p GoNews/data/news.csv -q "climate NOT hoax" -n 5
```

### Complex Queries
```bash
# Nested Boolean logic with parentheses
go run . -p GoNews/data/news.csv -q "(obama OR trump) AND policy" -n 10

# Multiple operators
go run . -p GoNews/data/news.csv -q "(climate OR environment) AND (policy NOT hoax)" -n 5
```

### Command-Line Flags

| Flag | Description | Default | Example |
|------|-------------|---------|---------|
| `-p` | Path to CSV file | `data/news.csv` | `-p GoNews/data/news.csv` |
| `-q` | Search query | `""` | `-q "climate change"` |
| `-n` | Max results to show | `10` | `-n 20` |
| `-stem` | Enable stemming | `false` | `-stem` |

### Example Commands

```powershell
# PowerShell examples
cd C:\Users\pc\Desktop\GoNews

# Simple search
go run . -p GoNews/data/news.csv -q 'election results' -n 5

# Political news
go run . -p GoNews/data/news.csv -q '(democrat OR republican) AND congress' -n 10

# Technology news
go run . -p GoNews/data/news.csv -q 'technology AND (AI OR artificial intelligence)' -n 15

# Exclude fake news
go run . -p GoNews/data/news.csv -q 'obama NOT fake' -n 10
```

## 🔧 Technical Architecture

### 1. Data Loading (`load.go`)
- Parses CSV files with format: `id,title,date,content`
- Creates `Document` structs for each article
- Handles malformed CSV entries gracefully

### 2. Text Analysis (`analyze.go`)
- **Tokenization**: Splits text into words using regex
- **Normalization**: Converts to lowercase, removes punctuation
- **Stemming** (optional): Porter stemmer algorithm for root words
- **Stopword Removal**: Filters common words (the, a, an, etc.)

### 3. Inverted Index (`index.go`)
- **Structure**: `map[term]map[docID][]positions`
- **Indexing**: O(n) time complexity for n documents
- **Search**: O(1) lookup for terms
- **Scoring**: TF-IDF algorithm for relevance ranking

#### TF-IDF Formula
```
score = (tf / doc_length) × log(N / df)

where:
  tf = term frequency in document
  doc_length = total tokens in document
  N = total number of documents
  df = document frequency (docs containing term)
```

### 4. Query Processing (`query.go`)
- **Parser**: Converts query to RPN (Reverse Polish Notation)
- **Operators**: AND, OR, NOT with proper precedence
- **Phrases**: Handles quoted strings as single units
- **Evaluation**: Stack-based RPN evaluation
- **Snippets**: Generates context previews around matched terms

### 5. Search Pipeline
```
Query Input
    ↓
Parse to RPN
    ↓
Evaluate Boolean Logic
    ↓
Get Matching Documents
    ↓
Calculate TF-IDF Scores
    ↓
Sort by Relevance
    ↓
Generate Snippets
    ↓
Return Top N Results
```

## 📊 Performance Metrics

Based on 6,335 news articles:

| Operation | Time | Details |
|-----------|------|---------|
| **Loading** | ~70ms | Parse CSV to memory |
| **Indexing** | ~3s | Build inverted index |
| **Search** | ~1ms | Query execution |
| **Total Startup** | ~3.1s | Load + Index |

**Memory Usage**: ~150MB for full index

## 📝 Dataset Format

### CSV Structure
```csv
id,title,date,content
0,"Article Title","2023-01-15","Full article text..."
1,"Another Article","2023-01-16","Article content..."
```

### Required Columns
- `id`: Unique identifier (integer)
- `title`: Article headline (string)
- `date`: Publication date (YYYY-MM-DD)
- `content`: Full article text (string)

## 🔍 Query Syntax Guide

### Basic Syntax
- **Single Term**: `climate`
- **Multiple Terms** (OR): `climate change`
- **Phrase**: `"climate change"`

### Boolean Operators
- **AND**: Both terms required → `climate AND policy`
- **OR**: Either term → `climate OR environment`
- **NOT**: Exclude term → `climate NOT hoax`

### Operator Precedence
1. NOT (highest)
2. AND
3. OR (lowest)

### Grouping with Parentheses
```
(climate OR environment) AND policy
obama AND (healthcare OR economy)
(trump OR biden) AND election NOT fraud
```

## 🧪 Example Queries & Results

### 1. Climate Change
```bash
go run . -p GoNews/data/news.csv -q "climate change" -n 3
```
**Output:**
```
Loaded 6335 docs in 70ms
Indexed 6335 docs in 3.2s
Search completed in 1ms — 1357 results

[2019-05-07] SHAMEFUL: Weather Channel Using Children... (score: 0.2284)
[2019-01-20] GOP avoids showdown over EPA climate... (score: 0.2101)
[2019-09-01] Obama Spends Earth Day In The Everglades... (score: 0.1538)
```

### 2. Political Query
```bash
go run . -p GoNews/data/news.csv -q "(obama OR trump) AND congress" -n 5
```

### 3. Technology Search
```bash
go run . -p GoNews/data/news.csv -q "technology AND innovation" -n 10
```

## 🚧 Future Enhancements

- [ ] Web UI with REST API
- [ ] Real-time indexing for new documents
- [ ] Fuzzy search / typo tolerance
- [ ] Multi-field search (title-only, content-only)
- [ ] Date range filtering
- [ ] Faceted search / filtering by category
- [ ] Synonym expansion
- [ ] Query suggestions / autocomplete
- [ ] Distributed indexing for larger datasets
- [ ] Persistent index storage (avoid reindexing)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📄 License

This project is open source and available under the MIT License.

## 👤 Author

**Raihan Kabir**  
Email: kabirraihan249@gmail.com

Built as a demonstration of information retrieval concepts and Go programming.

## 🙏 Acknowledgments

- Dataset: Kaggle Fake and Real News Dataset
- Inspired by classic IR systems like Lucene and Elasticsearch
- TF-IDF scoring based on standard IR textbooks

## 📧 Support

For issues, questions, or suggestions, please open an issue in the repository.

---

**Happy Searching! 🔍**
