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
              "type": "string"
            },
            "mime": {
              "type": "string"
            },
            "name": {
              "type": "string"
            },
            "path": {
              "type": "string"
            },
            "sha1": {
              "type": "string"
            },
            "sha256": {
              "type": "string"
            },
            "sha512": {
              "type": "string"
            },
            "size": {
              "type": "string"
            }
          }
        },
        "plugins": {
          "properties": {
            "archive": {
              "type": "object"
            },
            "av": {
              "type": "object"
            },
            "document": {
              "type": "object"
            },
            "exe": {
              "type": "object"
            },
            "intel": {
              "type": "object",
              "properties": {
                "virustotal": {
                  "dynamic": false,
                  "type": "object"
                }
              }
            },
            "metadata": {
              "type": "object"
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
//                 "type": "string"
//               },
//               "mime": {
//                 "type": "string"
//               },
//               "name": {
//                 "type": "string"
//               },
//               "path": {
//                 "type": "string"
//               },
//               "sha1": {
//                 "type": "string"
//               },
//               "sha256": {
//                 "type": "string"
//               },
//               "sha512": {
//                 "type": "string"
//               },
//               "size": {
//                 "type": "string"
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
