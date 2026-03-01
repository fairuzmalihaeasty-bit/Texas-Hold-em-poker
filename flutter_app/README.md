# Texas Hold'em Flutter Client

This is a minimal Flutter client that talks to the Texas Hold'em Go REST API included in the repository. It provides screens for hand evaluation, hand comparison, and Monte Carlo simulation.

Prerequisites
- Flutter SDK (stable) installed: https://flutter.dev/docs/get-started/install
- An emulator or device (Android emulator recommended for local API calls)

Quick start

1. Open a terminal and change directory to the Flutter app:

```bash
cd flutter_app
flutter pub get
```

2. Run on an Android emulator or connected device:

```bash
flutter run
```

Notes about the backend URL
- The app's API client uses the `BASE_URL` value in `lib/config.dart`. By default it's set to `http://10.0.2.2:8080` (Android emulator -> host machine).
- For a physical device or iOS simulator, update `lib/config.dart` to the appropriate host (e.g. `http://192.168.1.42:8080` or `http://localhost:8080` for iOS simulator).

Edit the client base URL

Open `lib/config.dart` and update the `BASE_URL` constant, for example:

```dart
const String BASE_URL = 'https://my-deployed-backend.example.com';
```

Build release APK (Android)

```bash
flutter build apk --release
adb install build/app/outputs/flutter-apk/app-release.apk
```

Troubleshooting
- If network requests fail on Android, ensure the emulator can reach the host and the server is running.
- On iOS, use `localhost` for a simulator; physical iOS devices must reach the host machine via LAN IP.
- If you change the API shape, update the screens in `lib/screens/` accordingly.

Development notes
- UI is intentionally minimal and uses text inputs for card entry. You can replace these with pickers or dropdowns in `lib/screens/`.
- The app depends on `http` package (see `pubspec.yaml`).

Questions or changes
- If you want a packaged example with card pickers or integrated auth, tell me which features to add and I can extend the app.
