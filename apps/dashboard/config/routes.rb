Rails.application.routes.draw do
  # Note: Devise deprecation warnings about hash arguments are from Devise internals (v4.9.4)
  # and will be resolved in a future Devise release for Rails 8.2 compatibility.
  # These warnings don't affect functionality.

  # Devise authentication
  devise_for :users, controllers: {
    registrations: "users/registrations",
    sessions: "users/sessions"
  }

  # Landing page
  root "home#index"

  # Health check
  get "up" => "rails/health#show", as: :rails_health_check

  # User dashboard namespace
  namespace :dashboard do
    root "overview#index"

    resources :api_keys do
      member do
        post :regenerate
        delete :revoke
      end
    end

    resource :usage, only: [ :show ], controller: "usage" do
      collection do
        get :by_endpoint
        get :by_date
        get :export
      end
    end

    resource :billing, only: [ :show, :update ], controller: "billing"
    post "billing/checkout", to: "billing#checkout", as: :checkout_billing
    post "billing/portal", to: "billing#portal", as: :portal_billing
    delete "billing/cancel_subscription", to: "billing#cancel_subscription", as: :cancel_subscription_billing

    resources :invoices, only: [ :index, :show ]

    resource :settings, only: [ :show, :update ] do
      member do
        delete :account # Delete account
      end
    end
  end

  # Admin panel (requires admin authentication)
  namespace :admin do
    root "dashboard#index"

    # Only allow authenticated admin users
    authenticate :user, ->(u) { u.admin? } do
      resources :users do
        member do
          post :suspend
          post :unsuspend
          post :ban
          post :make_admin
          post :remove_admin
        end

        resources :credit_adjustments, only: [ :new, :create ]
      end

      resources :api_keys, only: [ :index, :show ] do
        member do
          delete :revoke
        end
      end

      resources :abuse_reports do
        member do
          post :resolve
          post :investigate
        end
      end

      # Analytics namespace
      namespace :analytics do
        get :usage
        get :revenue
        get :system_health
      end
    end
  end

  # Webhooks (unprotected, verified by signature)
  namespace :webhooks do
    post "lemonsqueezy", to: "lemonsqueezy#create"
    post "stripe", to: "stripe#create"
    post "cloudflare", to: "cloudflare#create" # Usage sync from Worker
  end

  # Public pages
  get "docs", to: "home#docs"
  get "pricing", to: "home#pricing"
  get "about", to: "home#about"
  get "team", to: "home#team"
  get "privacy", to: "home#privacy"
  get "terms", to: "home#terms"
  get "contact", to: "home#contact"
  get "api_reference", to: "home#api_reference"
  get "changelog", to: "home#changelog"

  # New static pages
  get "blog", to: "home#blog"
  get "status", to: "home#status"
  get "glossary", to: "home#glossary"
  get "error_codes", to: "home#error_codes"
  get "faq", to: "home#faq"

  # Form pages
  get "suggest_api", to: "suggestions#new"
  post "suggest_api", to: "suggestions#create"
  get "talk_to_sales", to: "sales_inquiries#new"
  post "talk_to_sales", to: "sales_inquiries#create"

  get "examples", to: "examples#index"
  resources :apis, only: [ :index, :show ]
  resources :examples, only: [ :show ]

  # API Playground Proxy
  post "api/proxy", to: "api_proxy#create"
end
