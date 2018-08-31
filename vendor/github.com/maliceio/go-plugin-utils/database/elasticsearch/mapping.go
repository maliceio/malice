package elasticsearch

const mapping = `{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  },
  "mappings": {
    "samples": {
      "properties": {
        "file": {
          "properties": {
            "md5": {
              "type": "keyword"
            },
            "mime": {
              "type": "keyword"
            },
            "name": {
              "type": "keyword"
            },
            "path": {
              "type": "text"
            },
            "sha1": {
              "type": "keyword"
            },
            "sha256": {
              "type": "keyword"
            },
            "sha512": {
              "type": "keyword"
            },
            "size": {
              "type": "keyword"
            }
          }
        },
        "plugins": {
          "properties": {
            "archive": {
              "properties": {}
            },
            "av": {
              "properties": {}
            },
            "document": {
              "properties": {}
            },
            "exe": {
              "properties": {}
            },
            "intel": {
              "properties": {
                "virustotal": {
                  "dynamic": false,
                  "properties": {}
                }
              }
            },
            "metadata": {
              "properties": {}
            }
          }
        },
        "scan_date": {
          "type": "date"
        }
      }
    }
  }
}`

// const mapping = `{
//   "settings": {
//     "number_of_shards": 1,
//     "number_of_replicas": 0
//   },
//   "malice": {
//     "mappings": {
//       "samples": {
//         "properties": {
//           "file": {
//             "properties": {
//               "md5": {
//                 "type": "text"
//               },
//               "mime": {
//                 "type": "text"
//               },
//               "name": {
//                 "type": "text"
//               },
//               "path": {
//                 "type": "text"
//               },
//               "sha1": {
//                 "type": "text"
//               },
//               "sha256": {
//                 "type": "text"
//               },
//               "sha512": {
//                 "type": "text"
//               },
//               "size": {
//                 "type": "text"
//               }
//             }
//           },
//           "plugins": {
//             "properties": {
//               "archive": {
//                 "type": "object"
//               },
//               "av": {
//                 "type": "object"
//               },
//               "document": {
//                 "type": "object"
//               },
//               "exe": {
//                 "type": "object"
//               },
//               "intel": {
//                 "type": "object"
//               },
//               "metadata": {
//                 "type": "object"
//               }
//             }
//           },
//           "scan_date": {
//             "type": "date",
//             "format": "strict_date_optional_time||epoch_millis"
//           }
//         }
//       }
//     }
//   }
// }`
